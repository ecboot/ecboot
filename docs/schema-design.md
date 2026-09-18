# ECBOOT 数据库 Schema 设计（MVP v1）

社交电商平台第一版数据库设计。目标规模：初期 1 万用户、日订单 500。
技术栈：Spring Boot 4 + MySQL 8.4 LTS + MyBatis Plus + Redis 8 + Flyway。

迁移脚本位于 `apps/api/ecboot-start/src/main/resources/db/migration/`（Flyway 默认位置，应用启动自动执行，零配置）。

## 文件清单

| 文件 | 领域 | 表 |
|---|---|---|
| `V1__user_domain.sql` | 用户 | `user`、`user_address` |
| `V2__product_domain.sql` | 商品 | `product_category`、`product_brand`、`product_spu`、`product_sku` |
| `V3__inventory_domain.sql` | 库存 | `inventory`、`inventory_log` |
| `V4__cart_domain.sql` | 购物车 | `cart_item` |
| `V5__freight_domain.sql` | 运费 | `freight_template`、`freight_rule` |
| `V6__order_domain.sql` | 订单 | `trade_order`、`trade_order_item`、`trade_order_log` |
| `V7__payment_domain.sql` | 支付 | `pay_order`、`pay_callback_log` |
| `V8__after_sale_domain.sql` | 售后 | `after_sale_order` |
| `V9__promotion_domain.sql` | 促销 | `coupon`、`user_coupon` |
| `V10__admin_domain.sql` | 后台 | `admin_user`、`admin_role`、`admin_user_role`、`admin_permission`、`admin_role_permission`、`admin_login_log`、`admin_operation_log` |

共 26 张表。

### 增量迁移总览（V11~V21，特性 002 社交电商扩展）

| 版本 | 域 | 表/改动 | 前提 |
|---|---|---|---|
| V11 | 隐私安全 | `user` 改造（phone 密文化 + phone_hash 唯一）+ `user_login_log` | V1；存量数据环境须先跑应用侧密文化任务 |
| V12 | 评价 | `product_review` | V6 |
| V13 | 收藏足迹 | `user_favorite`、`user_footprint` | V2 |
| V14 | 通知 | `notify_task`、`user_message` | V1 |
| V15 | 物流 | `logistics_company` | — |
| V16 | 分销 | `user_relation`、`distribution_user`、`commission_rule`、`commission_record`、`user_account`、`account_log`、`withdraw_order`、`invite_record` | V2/V6 |
| V17 | 拼团 | `group_buy_activity`、`group_buy_team`、`group_buy_team_member` + `trade_order` 增列 | V6 |
| V18 | 秒杀 | `flash_sale_activity`、`flash_sale_item` | V2 |
| V19 | 积分等级 | `point_account`、`point_log`、`user_level_rule` + `user`/`trade_order`/`trade_order_item` 增列 | V1/V6 |
| V20 | 满减 | `promotion_activity`、`promotion_activity_ladder`、`promotion_activity_scope` + 订单增列 | V2/V6/V9 |
| V21 | 风控 | `risk_rule`、`risk_record` | V1 |
| V22 | 评审修复 | 满减全场行唯一性加固（生成列归一键）+ 关联列索引补齐 | V20 |
| V23 | 业务评审修复 | 地址区划码 / 售后数量 / 秒杀关联 / 包邮例外 / 关系链锁定 / 规格组合唯一 / 枚举扩位 | V22 |
| V24 | 拼团定价下沉 | 新表 `group_buy_item`（SKU 级成团价，对齐 flash_sale 模式），移除活动级 `group_price` | V17 |
| V25 | 余额消费 | 订单头/项 `account_amount` 抵扣列 + 账户流水枚举扩位（资金语义=用户资产消耗，非优惠） | V23 |
| V26 | 运营与配置 | 地区限售（SPU 黑名单 JSON）/ 积分滚动过期 / `notify_template` / `operation_banner`+`operation_floor` | V19/V2 |
| V27 | 休眠账户风控 | `user.last_login_at/last_active_at`（休眠分级依据）+ `risk_rule` 枚举扩位（5 休眠账户分级） | V1/V21 |
| V28 | 分享归因 | `share_record` 表 + `trade_order` 归因双列与卖家备注 + `user.share_code` + `invite_record.reward_trigger` | V1/V6 |

执行契约与唯一性清单见 `specs/002-social-commerce-expansion/contracts/schema-contracts.md`。

**域设计要点**（详细字段与状态机见 `specs/002-social-commerce-expansion/data-model.md`）：

- **隐私安全（V11）**：`user.phone` 密文 + `phone_hash` 唯一盲索引，仅完整手机号精确检索（个保法最小化）；登录日志只追加。
- **评价（V12）**：`UNIQUE(order_item_id)` 一项一评；商品快照列（spu_name/sku_specs）让评价不随商品软删漂移。
- **收藏/足迹（V13）**：`(user_id, spu_id)` 双唯一；足迹 90 天物理清理（`idx(last_view_at)`）。
- **通知（V14）**：一渠道一行 + 重试上限终态；站内信兜底必达。
- **分销（V16）**：关系链仅 `inviter_id` 单列两级封顶（ADR-0003）；佣金基数=订单项实付；账户可负；提现渠道单号唯一幂等。
- **拼团（V17）/秒杀（V18）**：完全复用订单/库存模型——团价/秒杀价走订单项快照；秒杀活动库存分账（`stock_count/sold_count` + CHECK），取消回补活动侧。
- **积分/满减（V19/V20）**：双账本（积分可负消耗、成长值只增）；订单优惠三构成恒等式 `promotion_amount ≡ coupon_amount + full_reduction_amount + point_amount`（头/项两层成立，尾差记末行）；满减范围用关系表支撑热路径命中查询。
- **风控（V21）**：规则与事件分离，事件多态关联业务对象、申诉四态；**站内搜索零建表**（商品既有结构足够）。

表总数：**54**（V1~V10 基线 26 + V11~V21 新增 28）；`user`/`trade_order`/`trade_order_item` 为增列改造。

### 注销匿名化规则（FR-003 / 合规红线 2）

注销 = 校验无进行中订单/售后/提现、佣金余额处理完毕 → 二次确认 → 匿名化落库：

| 列 | 处理 |
|---|---|
| `nickname` / `avatar` | 置 `'已注销用户'` / `''` |
| `phone` | 密文替换为占位密文（不再参与登录检索） |
| `phone_hash` | 墓碑值 `SHA256('deleted:' + user_id)`（确定性、全局唯一、不可反推原号）→ **原手机号可重新注册** |
| `password_hash` / `wx_openid` | 置 NULL（`uk_wx_openid` 唯一索引允许多 NULL） |
| 交易数据（订单/佣金/流水/日志） | **一律保留**（单据不可删，ADR-0002） |

墓碑唯一性由 `uk_phone_hash` 承载：反复注销/注册产生含 user_id 的不同墓碑，永不冲突。与 ADR-0002"严格唯一"的关系见该 ADR 的范围注记（普通软删≠注销）。

### V11 在存量数据环境的执行顺序（runbook，评审 I1）

V11 假定 `user` 为空表或已密文化（脚手架期即空库，已验证）。任何含明文手机号存量的环境**必须**按序执行：

1. 应用侧密文化任务：逐行生成 `phone` 密文与加盐 `phone_hash`（应用管理的影子流程，不在迁移链内）；
2. 前置自检通过：库内不存在未密文化的 `user` 行；
3. 执行 V11——其 `ADD COLUMN phone_hash NOT NULL`（无默认值）在未完成步骤 1-2 的非空表上会以隐式空串回填并在建唯一索引时**失败**，这是有意的防呆：迁移失败 = 前置未满足，环境被拒绝而非被污染。

（如未来需要无损形态，另起影子列三段式迁移；当前无此需求。）

### 评审修复与边界声明（V22）

- **恒等式仅约束增量订单**：V19/V20 之前创建的订单（`promotion_amount>0` 而三构成列=0）天然违反恒等式；数据巡检须按 V19 应用时间过滤增量。
- **拼团名额释放编排**：成员订单在待付款阶段取消（用户/超时）→ **物理删除团员行 + `member_count` 同事务递减**；留痕由 `trade_order_log` 承载，参团历史可经订单反查。支付后取消属售后域，不释放名额。
- **收藏复活语义**：取消收藏（软删）后再次收藏 = 复活原行（`deleted` 置 0），援引 ADR-0002 严格唯一，不产生第二行。
- **佣金防重复计佣**：由应用层承担（冲销负记录需要同键多行，数据库唯一键不可表达——契约 R5）。
- **关系链校验义务**：两级封顶是结构强制的"层级"约束；自邀（user_id=inviter_id）与 A↔B 互环需应用层拒绝（契约 R6）。
- V22 补齐关联列索引：`trade_order.promotion_activity_id`、`product_review.order_no`+`sku_id`、`user_message.biz_no`、`group_buy_team.leader_user_id`、`commission_record` 冲销双列、`user_footprint.spu_id`。

### 业务评审修订（V23，PM+技术双视角评审）

- **P0-① 地址区划代码**：`user_address` 增 `province_code/city_code/district_code CHAR(6)`（名称=展示快照、代码=匹配口径）——修复运费 `region_codes` 按代码匹配而地址只有名称的断裂。
- **P0-② 售后数量**：`after_sale_order.quantity`（默认 1）——买 3 退 1 的退款摊算/库存回补/佣金冲销依据。
- **P0-③ 秒杀关联**：`trade_order_item.flash_sale_item_id`（可空+索引）——限购校验与秒杀订单识别的数据落点。
- **P1-④ 包邮区域例外**：`freight_template.free_exclude_codes JSON`（不参与满额包邮的省份）。
- **P1-⑤ 关系链锁定**：`user_relation.locked/lock_time`——绑定窗口（保护期可换绑、期满锁定）的承载；换绑=UPDATE 本行，`UNIQUE(user_id)` 不受影响。
- **P1-⑦ 规格组合唯一**：`product_sku.specs_hash CHAR(64) AS (SHA2(specs,256)) STORED` + `UNIQUE(spu_id, specs_hash)`。**实证（2026-09-18）：MySQL JSON 存储层先规范化键序再参与生成列计算——键序不同的同组合同样被 1062 拒绝，属语义级唯一**（强于 V23 文件注释所述"依赖应用序列化稳定"，以本条为准）。
- 附：`point_log.biz_type` 枚举扩位（6 评价获得 / 7 注册赠送 / 8 邀请奖励）。
- **评审遗留（全部已决策，2026-09-18）**：~~拼团定价粒度~~→V24 SKU 级；~~佣金余额消费~~→V25；~~账号合并~~→登录时手机号优先归并（零 DDL，规则见 `ecboot-service-user/README.md`，不做迁移式合并）；~~地区限售~~→V26 SPU 黑名单；~~积分有效期~~→V26 滚动 12 个月（批次制为升级路径，触发条件：积分商城/兑换比例变化/财务审计要求）；~~通知模板~~→V26 `notify_template`（平台模板 ID + 参数契约建模）；~~运营位~~→V26 轻量两表。

### 分享归因（V28，社交电商第二轮评审 P0）

**归因模型（佣金计提三级判定，回答"这笔订单是谁带来的"）**：

| 优先级 | 类型 | 规则 |
|---|---|---|
| 1 | 分享归因 | 订单来自 `share_record` 分享触点（默认窗口 **7 天**，应用配置化）→ 归因人 = 分享人，**谁分享谁受益** |
| 2 | 关系链兜底 | 无归因时按 `user_relation` 两级计佣 |
| 3 | 自然流量 | 均无 → 不计佣 |

- `trade_order.attributed_user_id` + `attribution_type`（1分享/2关系链/3自然）落订单；`user.share_code`（唯一）承载海报/口令；`share_record` 只追加，是归因判定入口与分享行为分析源。
- **自购不算自己的归因**（防自刷）；自购返佣为独立计提规则，产品决策挂账（推广员自购按一级返本人？）。
- 首单激励：`invite_record.reward_trigger`（1 注册即发 / 2 首单后发，防刷主流形态）。
- 附带：`trade_order.seller_remark`（客服/仓库内部备注，买家不可见）。

### 休眠账户风控（V27，产品决策落地）

| 休眠时长（以 `last_active_at` 为准） | 动作 |
|---|---|
| ≥90 天 | 登录强制短信验证码（不再静默放行） |
| ≥180 天 | 休眠标记：敏感操作（改绑/提现/余额消费）二次核身 + 资金限频 |
| ≥365 天 | 深度休眠：提现/余额动用前重新实名核身；清理候选池 |
| ≥2 年 | 匿名化候选——**有余额/未结算佣金/未完结订单的账号绝不可清理**（资金安全 > 合规处置） |

- 判定口径用 `last_active_at` 而非登录态：小程序静默登录会掩盖"未登录"。
- 拦截原则：时长只定默认核身等级，**久未登录后的异常行为组合**（新设备/异地/密码错后成功/登录即改绑提现）才触发升级；异地登录提醒走 `notify_task`。
- 阈值与动作建议以 `risk_rule`（rule_type=5）配置化承载，不硬编码。

### 余额消费（V25，产品决策落地）

- **资金语义**：余额抵扣是**用户资产消耗**，与营销优惠（券/满减/积分）分账——既有恒等式全部保持，新增第四恒等式：
  **`pay_order.amount ≡ trade_order.pay_amount − trade_order.account_amount`**（现金实付 = 用户应付 − 余额抵扣，渠道单只收现金部分）。
- 订单头/项两层 `account_amount` 分摊（行分摊 = 售后按行原路退回依据，尾差记末行沿用既有规则）。
- 资金流程与提现同构（条件更新）：下单冻结（balance→frozen）→ 支付成功核销 → 取消/超时解冻回退 → 售后按行退回；`account_log.biz_type` 扩 6 消费冻结 / 7 消费完成 / 8 消费退回。
- 规则：仅可用余额（`balance − frozen` 正值部分）可消费，欠款不可（`user_account.balance` 注释已明示）。

## 全局约定

- **引擎/字符集**：InnoDB，`utf8mb4` / `utf8mb4_0900_ai_ci`。
- **主键**：`id BIGINT UNSIGNED AUTO_INCREMENT`（库存表例外：以 `sku_id` 为主键）。业务单号（`order_no`/`pay_no`/`after_sale_no`）独立生成、加唯一约束，为分库分表预留路由空间。
- **命名**：`领域前缀_snake_case`（如 `trade_order`、`pay_order`），规避 `order`/`user` 等关键字歧义。
- **时间**：`created_at` / `updated_at` 由数据库默认值维护；日志类表只有 `created_at`，只追加不修改。
- **软删除**：只有主数据（用户、地址、商品、分类、品牌、运费模板、优惠券、后台账号/角色/权限，以及 V11~V21 新增的物流公司、推广员、佣金规则、拼团活动、秒杀活动、等级规则、满减活动、风控规则、收藏）有 `deleted TINYINT(0/1)`（MyBatis Plus `@TableLogic`）；`product_review` 带软删属例外定位（用户内容，用户可删自己的评价）。**交易单据与审计日志永不删除**——取消是状态迁移，不是删除。
- **外键**：不建物理外键（分库演进友好、减少写入开销与死锁面），关联完整性由服务层事务 + 索引保证；关联列一律建索引。
- **金额与币种**：金额全部 `DECIMAL(10,2)`；币种 `currency CHAR(3)`（ISO 4217），落在资金头表（`trade_order`/`pay_order`/`after_sale_order`），明细继承头表；**优惠券仅 CNY**（多币种抵扣涉及汇率折算，明确不支持）。**禁止 float/double 存金额；应用层统一 `BigDecimal`，比较用 `compareTo`，禁止 `equals`/`==`**（已写入根 `AGENTS.md`，对所有 AI 会话生效）。
- **枚举**：`TINYINT` + 注释穷举取值，应用层常量/枚举类映射（MP `@EnumValue`）。
- **JSON 字段**：只用于"整存整取、不参与条件查询"的数据（商品图集、SKU 规格快照、回调原文、省份代码列表）；凡是要统计、分摊、按条件查询的（订单项、库存流水）必须拆表。

## 核心设计决策

### 1. SPU / SKU 拆分与规格存储

- `product_spu`：品类、品牌、图文详情、**规格定义** `spec_definitions JSON`（如 `[{"name":"颜色","values":["黑","白"]}]`）——声明规格轴与可选值。
- `product_sku`：一个规格组合一个 SKU，`specs JSON` 存组合快照 + `sku_no` 业务编码。价格/上下架挂在 SKU。
- 上下架两级：SPU `status`（整体）+ SKU `status`（单品），可售 = SPU 上架 AND SKU 启用 AND 库存 > 0。
- 不建 `sku_spec` 关系表：MVP 无"按规格值筛选"需求；将来做规格筛选时上 Elasticsearch（栈已就绪），不影响现有写入路径。

### 2. 订单项快照

`trade_order_item` 冻结下单时刻的商品事实：`spu_name`、`sku_name`、`sku_specs`、`sku_image`、`price`、`original_price`、分摊后的 `promotion_amount`/`pay_amount`。理由：商品改名/改价/下架/删除后，订单必须仍是自洽的历史凭证（客服对账、售后金额依据）。收货信息同理快照在 `trade_order` 上；运费快照 `freight_amount` + `freight_template_id`。

### 3. 库存三段式与防超卖（下单锁库存制）

字段：`total`（总库存=可售+锁定）、`locked`（下单未支付锁定），**可售 = total - locked** 为推导值不落库。数据库级约束 `CHECK (locked <= total)`。

**扣减时机：下单锁定、支付核销、取消/超时释放**（预扣制，杜绝支付成功却无货）：

```sql
-- 下单（本地事务内）：可售不足则失败
UPDATE inventory SET locked = locked + #{qty}
WHERE sku_id = #{skuId} AND total - locked >= #{qty};   -- affected=0 → 抛库存不足
-- 支付成功：核销
UPDATE inventory SET total = total - #{qty}, locked = locked - #{qty}
WHERE sku_id = #{skuId} AND locked >= #{qty};
-- 取消/超时：释放
UPDATE inventory SET locked = locked - #{qty} WHERE sku_id = #{skuId} AND locked >= #{qty};
-- 退货退款验货通过：回补
UPDATE inventory SET total = total + #{qty} WHERE sku_id = #{skuId};
```

全部原子条件更新（单行行锁串行化，无乐观锁重试复杂度），每次变更同事务写 `inventory_log`（前后值快照）供对账。

#### 热点 SKU 行锁竞争解决方案（Redis 预扣挡板）

当前量级（日单 500）单行行锁竞争可忽略，DB 条件更新即是完整方案。当单 SKU 并发下单达到热点（如秒杀/爆品抢购）时，**不改表结构**启用 Redis 预扣：

```
key:  stock:avail:{skuId}  = 可售数（与 DB total-locked 对齐，事件+定时刷新）

预扣 Lua（原子，Redis 8 / EVAL）:
  local v = tonumber(redis.call('GET', KEYS[1]))
  if not v or v < tonumber(ARGV[1]) then return -1 end
  return redis.call('DECRBY', KEYS[1], ARGV[1])

流程: Redis预扣成功 → DB事务(锁库存+建订单+锁券)
      → DB失败/回滚 → Redis 补偿 INCRBY 回加（对账任务兜底漏网）
取消/超时释放: Redis INCRBY 回加该数量
支付核销: 不动 Redis（可售数在锁定时已减）
```

要点：**Redis 只挡流量，DB 条件更新仍是唯一事实源**；任何时刻 DB 才是真相，Redis 未命中或数值异常时回退纯 DB 路径。对账任务定时校验 `Redis(avail) ≈ DB(total-locked)`，偏差告警并重置。再升级（秒杀级）做库存分桶（`inventory` 加 bucket 维度拆行摊薄单行热度），属触发式演进。

### 4. 运费模板

- `freight_template`：计费方式（1 按件数 / 2 按重量）+ `free_threshold` 满额包邮（NULL=不包邮）。
- `freight_rule`：区域差异化规则，`region_codes JSON` 为适用省级代码列表（空数组=全国兜底，每模板至多一条，应用层保证）；首段（first_unit/first_fee）+ 续段（continue_unit/continue_fee）。
- 商品关联：`product_spu.freight_template_id`（NULL=包邮）。
- 计算规则（应用层）：①商品金额 ≥ free_threshold → 0；②按 charge_type 折算计费量 Q（按重时 `sku.weight` 缺失按 1 件折算）；③匹配收货省规则（无匹配用兜底）；④`运费 = first_fee + Q≤first_unit ? 0 : CEIL((Q-first_unit)/continue_unit)*continue_fee`。下单时算出并快照进订单。

### 5. 幂等总则（五个入口全覆盖）

| 入口 | 机制 |
|---|---|
| **下单** | 确认页领 `request_token`（Redis `SET NX`），提交时 `DEL` 抢占；DB `UNIQUE(user_id, request_token)` 兜底（NULL 不参与唯一，历史数据零影响） |
| **支付回调** | ①条件更新抢占 `UPDATE pay_order SET status=20... WHERE pay_no=? AND status=10`，affected=0 即重复通知直接 ACK；②`UNIQUE(pay_channel, channel_trade_no)` 防换号重放；③金额+币种校验（回调金额=amount=订单 pay_amount，currency 一致，不符告警人工）；④原文全量落 `pay_callback_log` 可重放 |
| **退款回调** | 渠道幂等键 `out_refund_no = after_sale_no`（重试不产生新退款）；状态机条件更新 `WHERE status=40`，affected=0 即重复直接 ACK |
| **领券** | `UPDATE coupon SET received_count=received_count+1 WHERE id=? AND (total_count=0 OR received_count<total_count)`，affected=0 即已领完/超领 |
| **定时任务**（超时取消/自动收货） | 全部条件更新（`WHERE status=10 AND created_at<...`），天然幂等可重跑 |

### 6. RBAC 与审计

- **RBAC 五表**：`admin_user`（`is_super` 跳过校验）↔ `admin_user_role` ↔ `admin_role` ↔ `admin_role_permission` ↔ `admin_permission`（菜单/按钮/接口统一权限树，`code` 如 `product:spu:create`，接口权限与 API 路径对应）。
- **登录审计** `admin_login_log`：成功/失败全记（含账号不存在的用户名尝试，防暴力破解分析），同步写（登录低频）。
- **操作审计** `admin_operation_log`：AOP 切面拦截写操作，**异步写**（队列/线程池），写失败不阻断业务；`username` 冗余快照，账号删除后日志仍可读；`request_params` 脱敏后留档。
- 初始超级管理员由应用启动初始化（密码哈希须应用生成，不在 SQL 中硬编码）。

### 7. 优惠券规则

下单成功：`user_coupon` → 已使用并绑定 `order_no`（与订单同事务）；订单取消/超时：→ 已退回（可再用，过期时间不变）；优惠分摊按行商品金额比例，尾差记最后一行。

## 状态机

### 订单 `trade_order.status`

| 迁移 | 触发 | 副作用 |
|---|---|---|
| 10 待付款 → 20 待发货 | 支付成功回调 | 核销库存、记 `pay_time` |
| 10 → 90 已取消 | 用户 / 30 分钟超时(Quartz) / 管理员 | 释放库存、退回优惠券 |
| 20 → 30 待收货 | 管理员发货 | 写物流公司/单号 |
| 30 → 40 已完成 | 用户确认收货 / 发货后 7 天自动 | 记 `finish_time` |

20/30/40 均可发起售后：**不改主状态**，由 `refund_status`（0 无 / 1 部分退款 / 2 全额退款）+ 售后单表达。所有迁移用条件 UPDATE，非法迁移 affected=0 即拒绝。

### 支付单 `pay_order.status`

10 待支付 → 20 成功（回调验签通过，终态）；10 → 30 失败（可重新发起新支付单，旧单关闭）；10 → 90 已关闭（订单取消或过期）。同一订单同时至多一个"待支付"支付单（服务层保证）。

### 售后单 `after_sale_order.status`

10 待审核 → 同意：仅退款直接进 30；退货退款进 20 待买家寄回 → 确认收货 → 30 待退款 → 渠道退款受理 → 40 退款中 → 退款成功回调 → 50 已完成（回补库存、更新订单 `refund_status`）。10 → 90 已拒绝（终态，可重新申请）；未终态可 → 91 已撤销（买家主动 / 寄回超时）。

## 事务边界

**强一致（同一本地事务）**：下单（订单+订单项+锁库存+核销优惠券）；支付回调（支付单+订单+库存核销）；售后完成（售后单+库存回补+订单 refund_status）；领券（coupon+user_coupon）；RBAC 授权变更。

**最终一致（事件/异步）**：`spu.sale_count` 销量冗余累加；超时取消/自动收货扫描（Quartz）；操作审计日志（AOP 异步）；对账任务（库存 Redis↔DB、回调日志重放）。

## 预留与演进

**已预留**：`sku.barcode`（扫码）；`after_sale_order.type=3 换货`；`coupon.type=3 折扣券` 及适用范围；`user.wx_unionid`（多端打通）；`inventory.warn_count`（库存预警）；`trade_order` 物流字段与 `order_channel`；运费按重计费（`sku.weight` 已备）。

**触发式演进**：商品搜索复杂化 → ES；单表过千万行（日单 500 约 5+ 年）→ 按 `user_id` 分片（业务单号分片友好）；秒杀热点 → Redis 预扣（方案见上）→ 库存分桶；多币种真实需求 → SKU 按币种定价属国际化范畴，需重构价格模型，明确不在本设计范围。
