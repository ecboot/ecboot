# Feature Specification: 营销后台（批次10）

**Feature Branch**: `016-marketing-admin`

**Created**: 2026-09-22

**Status**: Draft

**Input**: 批次 10：营销后台。**30 个 admin 端点**（当前全为 `CodeNotImplemented` 桩；service 接口 `ICouponLogic`/`IActivityLogic` 已有定义、全仓无实现者）：**券模板管理 6**（列表/详情/创建/修改/软删/领取记录）、**满减活动 5**（列表/详情/创建/修改/软删，档位+范围嵌套提交）、**拼团 5 / 秒杀 5 / 砍价 5**（各 CRUD + 场次商品全量替换设置）、**助力 4**（CRUD）。权限点已全部在 000032 RBAC 种子就位（`promotion:coupon:{read,create,update,delete}`、`promotion:{fullreduction,groupbuy,flashsale,bargain,assist}:manage`），**预期零权限迁移**。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 券模板管理（Priority: P1）🎯 MVP

运营者创建券模板（满减/无门槛两种类型；总量 0=不限、每人限领；有效期"固定区间"或"领取后 N 天"），可修改（名称/门槛/抵扣/总量/限领/启停）、软删、按状态筛选分页查看列表，并下钻查看某券的**领取/使用记录**（谁领的、未用/已用/过期/退回、核销单号）。C 端领券与计价（批次 06 已接 `ICouponQuery`）读同一张表——后台配置即刻生效。

**Why this priority**: 券是唯一有"记录下钻"与 C 端**已在线消费者**（下单核销计价）的管理对象——配错即影响线上交易，是本批价值最高、约束最硬的一块。

**Independent Test**: 创建券模板 → 列表可见（含已领数）→ 修改停发 → C 端首页"可领券"不再出现 → 记录页能看到既有领取行。

**Acceptance Scenarios**:

1. **Given** 合法创建参数（type=1 满减、threshold/discount、validType=1 固定区间），**When** 创建，**Then** 返回券 ID，列表可见且 `receivedCount=0`。
2. **Given** type=2（无门槛），**When** 创建不带 threshold，**Then** 成功；**When** type=1 不带 threshold，**Then** 拒绝（参数校验）。
3. **Given** validType=2（领取后 N 天），**When** 创建带 validDays，**Then** 成功；**When** validType=2 而 validDays≤0，**Then** 拒绝。
4. **Given** 既有券模板，**When** 修改 Status=0（停发），**Then** 生效——C 端"可领券"列表不再出现该券。
5. **Given** 已有会员领取过的券，**When** 软删，**Then** 券从后台与 C 端消失，既有 `user_coupon` 行不受影响。
6. **Given** 一张已被领取的券，**When** 查看记录，**Then** 分页返回领取行（用户/状态/领取时间/核销单号）。

---

### User Story 2 - 满减活动管理（档位+范围）（Priority: P1）

运营者创建满减活动：时间窗 + **多档位**（满 X 减 Y，同活动同门槛唯一）+ **适用范围**（全场 / 指定分类 / 指定商品，可多行；空=全场）。修改按**全量替换**档位与范围；详情返回完整档位与范围。C 端列表（批次 09）与下单计价共享此配置。

**Why this priority**: 满减是唯一**嵌套配置**（档位+范围）且**直接参与下单计价**（`calcFullReductionFen`）的活动——批次 09 复审移交的"配置面必须做对"即指此处（scope 行全链路写对，下单侧认范围属后续批次清偿但依赖本批把数据配对）。

**Independent Test**: 创建带 2 档 + 指定分类范围的活动 → 详情回读一致 → 重复门槛被拒（50008）→ 修改为全量替换后的档位/范围 → 详情一致。

**Acceptance Scenarios**:

1. **Given** 档位 [(100,10),(200,30)] + 范围 [全场]，**When** 创建，**Then** 详情回读档位按门槛升序、范围 1 行全场。
2. **Given** 同一活动重复提交门槛 100 两行，**Then** 拒绝（50008 档位门槛重复）。
3. **Given** 范围 [分类 C1]，**When** 创建时 targetId 为空或非数字，**Then** 拒绝（参数校验）；scopeType=1（全场）时 targetId 必须为空。
4. **Given** 既有活动，**When** 修改为新的档位/范围集合，**Then** 旧档位/范围被**全量替换**（不多不少）。
5. **Given** 进行中活动，**When** 软删，**Then** C 端满减列表不再出现。
6. **Given** 时间窗非法（start ≥ end），**When** 创建，**Then** 拒绝。

---

### User Story 3 - 拼团/秒杀/砍价活动管理（Priority: P1）

三类玩法活动同构管理：先建活动（拼团/砍价绑一个 SPU；秒杀纯时间窗），再用 `/{id}/items` **全量替换**场次商品（SKU 级：拼团=成团价；秒杀=秒杀价/限量/限购；砍价=起始价/底价/最大刀数）。列表按状态筛选分页。秒杀场次商品同 SKU 同活动唯一（`uk_activity_sku`）；活动分账库存的限量配置不得与商品库存语义混淆（只填"活动限量"，剩余=限量−已售）。

**Why this priority**: 三类同构一次收齐；秒杀场次商品是 C 端秒杀下单的**前置配置**（批次 09 已上线的消费端等着这批把配置面接通）。

**Independent Test**: 建秒杀活动 → 设置 2 个场次商品 → 列表/详情可见 → 重复同 SKU 设置被拒（唯一键）→ 全量替换为 1 个 → 旧场次商品消失。

**Acceptance Scenarios**:

1. **Given** 秒杀活动 + 场次商品 [(SKU1, 5.00, 10, 2)]，**When** 设置 items，**Then** `flash_sale_item` 两列落库正确（秒杀价/限量/限购）。
2. **Given** 同活动重复设置同一 SKU，**When** 在同一次提交中出现两行，**Then** 拒绝（唯一约束 → 业务码）。
3. **Given** 已设置 2 个场次商品，**When** 全量替换为 1 个，**Then** 被移除的场次商品行删除（若已被下单占用，按"进行中活动禁止替换/删除保护"处理——见 Edge Cases）。
4. **Given** 拼团活动，**When** 设置场次商品缺 groupPrice 或价格非法（>2 位小数/非正数），**Then** 拒绝。
5. **Given** 砍价场次商品 floor ≥ original 或 maxCutCount<1（非 0），**Then** 拒绝（0=不限由金额收敛保证）。
6. **Given** 拼团 GroupSize<2，**When** 创建，**Then** 拒绝（参数校验 min:2）。
7. **Given** 进行中的活动，**When** 软删，**Then** C 端对应列表立即消失（同表同口径）。

---

### User Story 4 - 助力活动管理（Priority: P2）

运营者创建助力活动：奖励类型（1=券：指定券 ID；2=积分：配置积分数存活动 config）、所需人数、每人可发起次数（**schema 事实：`uk_activity_user` 钉死"每人一次"，配置 >1 不生效**——契约必须声明该约束）、时间窗。修改与软删同构。

**Why this priority**: 结构最简（无场次商品），但承载批次 09 的"发奖意图"消费端（`IAssistReward` 端口已投递、实际发放属后续批次）——配置面就位后全链路语义闭合。

**Independent Test**: 创建 rewardType=2 的助力活动（积分数 100）→ 详情回读 rewardDesc 含积分数 → perLimit>1 的创建被接受但契约注释与列表展示明确"每人一次"语义。

**Acceptance Scenarios**:

1. **Given** rewardType=1 + rewardRef=券ID，**When** 创建，**Then** `assist_activity.reward_ref` 落该券 ID。
2. **Given** rewardType=2 + pointAmount=100，**When** 创建，**Then** `reward_ref=0` 且积分存入 `config` JSON（C 端 `RewardDesc` 可读出）。
3. **Given** rewardType=1 而 rewardRef 为空/非法，**Then** 拒绝。
4. **Given** requiredCount<1，**When** 创建，**Then** 拒绝（min:1）。
5. **Given** 修改 perLimit=3，**When** 更新成功，**Then** 契约层返回的 perLimit 如实为 3，但**实现注释与文档声明**">1 因唯一键不生效（每人一次）"。

---

### Edge Cases

- **时间窗**：所有活动 start ≥ end 一律拒绝；时间解析非 `YYYY-MM-DD HH:MM:SS` 拒绝。
- **删除保护**：场次商品已被订单引用（秒杀单 `trade_order_item.flash_sale_item_id` 指向它）时，活动软删照常（软删不停已生成订单的履约），但**进行中活动的场次商品替换**须保护已售 SKU（移除已售行 → 拒绝，防活动库存账悬空）。
- **券停发 vs 删除**：停发（status=0）可逆；软删（deleted=1）后同名单可重建，既有领取记录保留可查。
- **ID 越权/不存在**：修改/删除/详情/记录 对不存在 ID → 404 类业务码（50003/10006 对齐既有活动域口径），不泄露存在性差异。
- **per_limit>1（助力）**：接受配置但契约声明不生效（schema 唯一键兜底"每人一次"）——不静默改值、不报错。
- **分页越界**：page 超界返回空列表 + 正确 total。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-1（券）**: 券模板 CRUD + 状态筛选分页列表（含已领数）+ 详情 + 领取记录分页下钻；类型/有效期方式的条件校验（见 US1 场景）；停发即时影响 C 端可领口径。
- **FR-2（满减）**: 档位+范围嵌套创建、全量替换更新、详情回读（档位按门槛升序）；同活动门槛唯一（50008）；范围三型校验（全场无 target、分类/商品必须有合法 target）。
- **FR-3（场次商品）**: 拼团/秒杀/砍价的 `/{id}/items` 全量替换：同活动同 SKU 唯一（撞唯一键转业务码）；秒杀需 flashPrice+stockCount（perLimit 缺省 1）；拼团需 groupPrice；砍价需 original>floor 且 maxCutCount≥0；金额一律 `DECIMAL(10,2)` 口径、拒绝 >2 位小数。
- **FR-4（助力）**: 奖励两型落库（券→reward_ref；积分→config JSON）；perLimit 契约语义声明（>1 不生效）；CRUD 同构。
- **FR-5（删除）**: 全部为软删（deleted=1）；进行中活动的场次商品替换含"已售保护"。
- **FR-6（权限）**: 按路由注释挂权限点——券读（列表/详情/记录）`promotion:coupon:read`、券写 `promotion:coupon:{create,update,delete}`、其余四类活动全端点 `promotion:{fullreduction,groupbuy,flashsale,bargain,assist}:manage`；权限点已全部就位（000032），无需权限迁移。
- **FR-7（列表口径）**: 管理列表支持 Status 筛选（含全部）与分页；列表/详情/范围查询同源（避免口径漂移）；列表内档位/范围/SPU 简介尽量批量查询（收敛批次 09 记账的 N+1，同文件域内）。
- **FR-8（审计基线）**: 全部写操作走后台会话鉴权（operator 语义沿用 000029/000038 既有列——活动表无 operator 列，审计依赖 admin 会话与路由层，不新增列，与批次 03/04 管理面同口径）。

### Key Entities *(include if feature involves data)*

- `coupon`（券模板：类型/门槛/抵扣/总量/已领/限领/有效期两型/状态/软删）+ `user_coupon`（领取记录：状态/核销单号）
- `promotion_activity` + `promotion_activity_ladder`（档位，uk_activity_threshold）+ `promotion_activity_scope`（范围三型）
- `group_buy_activity` + `group_buy_item`（uk_activity_sku；成团价）+ `group_buy_team`（只读引用，成团/解散属后续批次）
- `flash_sale_activity` + `flash_sale_item`（uk_activity_sku + chk_sold_le_stock；秒杀价/限量/限购）
- `bargain_activity` + `bargain_item`（uk_activity_sku；起始价/底价/最大刀数）
- `assist_activity`（reward_type/reward_ref/required_count/per_limit/config）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-1**: `check-stub` 对账 **admin 渠道本批 30 端点桩数 58 → 28**（四渠道合计 69 → 39）。
- **SC-2**: 全量 `go test ./...` 两连跑全绿（含 01~09 与交易链路零退化）；`golangci-lint run` **0 issues**。
- **SC-3**: 每类活动至少一条"创建→配置→C 端可见性联动"的端到端测试（后台停发/软删 → C 端列表消失）。
- **SC-4**: 零新迁移、零新权限种子（勘察已确认就位）；生成物零变更或仅限真变更。

## Assumptions

- 拼团**成团/解散/退款**流程、秒杀/拼团的**参与记录管理**不在本批（13 批规划未含，属后续批次/待规划）。
- 下单侧满减"认范围"（批次 09 复审移交的 Important）属 `order_impl.go`（012 文件），本批只保证配置面全链路正确；计价侧清偿在台账挂账，见 PROGRESS §六 P1（13 批收官后归口独立横切批——原指针失效已勘误）。
- 助力**发奖实际发放**（积分/券入账）属批次 11 账户域，本批只做配置面。
- 管理面写操作的操作者审计沿用既有会话与路由层能力，不新增 operator 列（批次 03/04 同口径）。
