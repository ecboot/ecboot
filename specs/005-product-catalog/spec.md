# Feature Specification: 商品目录域实现（浏览端点 + 后台商品管理与库存）

**Feature Branch**: `（未创建）`

**Created**: 2026-09-21

**Status**: Draft

**Input**: User description: "商品目录域实现——shop 渠道浏览端点 + admin 商品管理端点的service/controller/仓储实现与测试"

## User Scenarios & Testing *(mandatory)*

> 首个完整业务域实现：把 004 契约与 005 服务层接口设计变成可运行的商品目录能力。
> 范围 = 商品浏览（C 端读）+ 后台商品管理（admin 写）+ 库存查询/调整；评价"写入"不在本特性（评价列表/汇总的"读"在内）。

### User Story 1 - 商品浏览（Priority: P1）

游客/会员可通过接口浏览三级分类树、品牌列表、商品列表（分类/品牌/价格区间筛选、销量/价格/上新排序）、商品详情（SKU 列表含规格/价格/可售态、规格定义、评价汇总）；被管理端下架的商品从列表消失、详情返回下架语义。

**Why this priority**: 转化漏斗入口，且为 admin 管理提供验收的"观察面"。

**Independent Test**: 空库起步：经后台接口创建分类/品牌/SPU/SKU 并上架后，浏览端点可完整还原目录结构并检索到该商品。

**Acceptance Scenarios**:

1. **Given** 后台已创建三级分类与启用品牌，**When** 请求分类树与品牌列表，**Then** 返回完整结构与排序一致。
2. **Given** 一个上架 SPU（2 个启用 SKU）,**When** 请求商品列表与详情，**Then** 列表含该商品（价格区间取 SKU 最小~最大），详情返回 2 个 SKU 及可售状态。
3. **Given** 商品详情被会员查看，**When** 请求，**Then** 浏览足迹被记录（重复查看不产生新行）。

---

### User Story 2 - 后台商品管理（Priority: P1）

管理员可以：管理三级分类（增改删，有商品引用禁删）；管理品牌；创建/编辑 SPU（规格定义/图文/运费模板/限售区域）；管理 SKU（新增含规格组合冲突校验、编辑、启停、删除软删）；两级上下架（SPU 上架要求至少一个启用 SKU；SKU 启停）。

**Why this priority**: 商品数据的唯一入口；浏览端点的数据由它产生。

**Independent Test**: 完成"建分类→建品牌→建 SPU→加 2 个 SKU→上架"的接口序列后，重复规格组合被拒绝（30006）、无启用 SKU 时上架被拒（30005）、下架后 C 端列表不再返回。

**Acceptance Scenarios**:

1. **Given** 同一 SPU 下已存在"黑/M"组合，**When** 再新增相同组合 SKU，**Then** 拒绝并返回 30006。
2. **Given** SPU 仅有一个禁用 SKU，**When** 执行上架，**Then** 拒绝返回 30005。
3. **Given** 分类下有在架商品，**When** 删除该分类，**Then** 拒绝（软删不执行）。

---

### User Story 3 - 库存查询与调整（Priority: P2）

管理员可查询 SKU 库存（总/锁定/可售）、按条件筛选、调整库存（有符号增量，留痕）、查看预警列表（可售≤阈值）。

**Why this priority**: 商品上架后即可售，库存是售卖的闸门；调整留痕是对账依据（ADR-0001）。

**Independent Test**: 调整 +100 后列表与详情可售数同步 +11 并有一条调整流水；负向超调被拒绝（30007）。

**Acceptance Scenarios**:

1. **Given** SKU 可售 10，**When** 调整 -20，**Then** 拒绝返回 30007。
2. **Given** 调整 +50 成功，**When** 查询库存列表与预警，**Then** 可售数更新为 60 且脱离预警；inventory_log 新增一条 operator=后台账号的记录。

---

### User Story 4 - 限售区域与搜索（Priority: P3）

管理员可设置 SPU 禁售省级代码列表；游客搜索接口支持关键词+筛选+排序（V1 以数据库条件查询实现契约语义，ES 为演进路径不改契约）。

**Why this priority**: 合规能力（限售）与体验增强（搜索）；依赖 US1/US2 的数据基础。

**Independent Test**: 设置某省禁售后，该省收货地址在下单校验（后续交易特性消费）与商品详情响应中可见限售标记；搜索按关键词命中并正确排序。

**Acceptance Scenarios**:

1. **Given** SPU 设置禁售代码 ["650000"]，**When** 查询详情，**Then** 响应含限售标记供前端置灰。
2. **Given** 关键词命中多个商品，**When** 按销量排序搜索，**Then** 返回列表按 sale_count 降序。

### Edge Cases

- 软删的 SPU/SKU/分类/品牌在 C 端全部不可见；在管理端默认也不出现（需显式筛选才可见的历史数据场景 V1 不做）。
- 价格区间筛选边界为闭区间（含 min 与 max）；价格以 SKU 最低价参与比较。
- 商品列表价格区间 = 启用 SKU 的价格 min~max；全部禁用时价格显示"—"且不可下单。
- 分类禁用后 C 端树不返回该节点；其下商品在列表中仍可见（分类筛选不可达）。
- 排序稳定性：同销量按 id 升序；分页参数越界返回空列表而非报错。
- 规格组合唯一（specs_hash）在并发新增同组合时的重复键错误映射为 30006。
- 库存调整并发：条件更新（total >= -delta 校验负向），留痕与更新同事务。
- 关键词搜索对特殊字符（%、_）转义；空关键词返回全部（分页）。

## Requirements *(mandatory)*

### Functional Requirements

**浏览（C 端公开）**

- **FR-001**: 分类树接口 MUST 返回启用分类组成的三级树（禁用与软删不可见），同级按 sort 升序。
- **FR-002**: 品牌列表 MUST 仅返回启用品牌。
- **FR-003**: 商品列表 MUST 仅返回上架 SPU，支持分类/品牌筛选、价格区间（闭区间）、排序（综合/销量/价格升/价格降/上新）与分页；每项含价格区间（启用 SKU min~max）。
- **FR-004**: 商品详情 MUST 返回 SPU 信息、全部启用 SKU（规格/价格/可售态=SPU上架 AND SKU启用 AND 可售>0）、规格定义、评价汇总（均值/分布/总数，仅审核通过评价）与限售标记；不暴露成本价。
- **FR-005**: 搜索接口 MUST 支持关键词（名称 LIKE，特殊字符转义）+ FR-003 同款筛选排序；V1 以数据库查询实现，契约不变（ES 演进）。
- **FR-006**: 会员查看详情 MUST 记录浏览足迹（UPSERT 语义）。

**后台管理（admin，写操作带权限点）**

- **FR-007**: 分类管理 MUST 提供树查询（含禁用）与增改删；删除时存在商品引用（含软删商品）MUST 拒绝（30004）。
- **FR-008**: 品牌管理 MUST 提供列表/增改删；名称唯一由数据库约束保证。
- **FR-009**: SPU 管理 MUST 提供列表（状态/分类/关键词筛选）、创建（生成 spu_no）、详情（含成本价与限售标记）、编辑、软删、上下架与限售设置；上架 MUST 校验至少一个启用 SKU（30005）。
- **FR-010**: SKU 管理 MUST 提供新增（同 SPU 规格组合唯一，冲突 30006；生成 sku_no）、编辑、启停、软删。
- **FR-011**: 库存 MUST 提供列表（可售推导）、调整（有符号增量 + 操作留痕同事务 + 负向超调拒绝 30007）、预警查询。
- **FR-012**: 限售区域 MUST 以 SPU 级省级代码黑名单承载（空=全国可售）。

**质量**

- **FR-013**: 全部接口遵循统一响应契约与既有错误码分段（3xxxx 商品域）。
- **FR-014**: 管理端写接口 MUST 校验权限点（product:category/brand/spu/sku:*）与操作留痕口径。
- **FR-015**: service 层方法 MUST 遵循红绿 TDD（先测试后实现），测试配置确定性注入。

### Key Entities *(include if feature involves data)*

- 既有实体（零新表）：product_category / product_brand / product_spu / product_sku / inventory / inventory_log / product_review（读）。
- 契约结构已定义于 `internal/service/shop`（ProductQuery/ProductCard/ProductDetail/SkuInput 等）与 `internal/model`。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 全部端点在统一响应契约下工作，错误码使用 3xxxx 段且与契约表一致。
- **SC-002**: "后台建数据 → C 端可见"的正反向场景全部通过（上架后可见、下架后不可见、软删后不可见）。
- **SC-003**: 规格组合唯一（30006）、无启用 SKU 上架（30005）、分类引用禁删（30004）、库存负向超调（30007）四条防御各有正反测试。
- **SC-004**: 商品列表价格区间、排序、筛选正确率 100%（构造数据集断言）。
- **SC-005**: 浏览详情产生足迹记录且重复浏览不增行（UPSERT 验证）。
- **SC-006**: service 层测试覆盖全部 IProductLogic/IInventoryLogic 方法（红绿交付），`make build && make test` 全绿。

## Assumptions

- 搜索 V1 以数据库条件查询实现契约语义；ES 为演进路径（契约不变）。
- 图片仅存 URL（OSS/上传组件另立特性）；测试用占位 URL。
- 评价汇总读 product_review（审核通过），评价写入属后续评价特性。
- 运费模板在本特性仅引用 ID（计算逻辑随交易特性实现）。
- 后台操作留痕依赖既有审计切面口径（operator 记录到 inventory_log；操作日志表由治理特性对接）。
- 库存初始值随 SKU 创建置 0，由后台调整录入（不提供创建时直接设库存，保证全量留痕）。
