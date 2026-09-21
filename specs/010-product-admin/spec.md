# Feature Specification: 商品目录后台与 C 端浏览收口（批次04）

**Feature Branch**: `（未创建）`

**Created**: 2026-09-21

**Status**: Draft

**Input**: User description: "010-product-admin 商品目录后台与 C 端浏览收口：admin 渠道 spu/sku/类目/品牌/库存共 22 端点 + shop 渠道 C 端浏览 5 端点的实现与连线（见 specs/PROGRESS.md 批次04）"

## User Scenarios & Testing *(mandatory)*

> **本批性质（勘察结论）**：服务层实现**大半已有**（005 商品目录特性遗产）——`ProductLogicImpl` 已实现
> 分类树/品牌/SPU/SKU 的管理方法（19 个）与 C 端浏览方法（5 个），但**从未被 controller 调用过**
> （端点全为桩）。本批因此是"**连线收口 + 库存补齐**"：把既有实现接到契约端点、补齐库存后台、
> 挂上权限；不重写既有业务逻辑。
>
> 全部行为基于既有契约：product_category/brand/spu/sku/inventory 表（000002/000003/000034）、
> service 接口（`IInventoryLogic` 等）、权限码 `product:*` 与 `inventory:*`（000032 种子）。
> 库存的**交易侧**操作（下单锁定/支付核销/取消释放/售后回补）由 006 交易链路以事务内条件更新实现，
> 不在本批范围；本批只做**后台管理面**（查询/预警/调整）。

### User Story 1 - 后台商品目录管理（Priority: P1）🎯 MVP

授权运营维护商品目录：分类（树形增删改）、品牌（列表/增改删）、SPU（列表/详情/增改删/上下架/限售区域）、SKU（增改删/启停）。服务实现已在 005 交付并有测试，本批接入契约端点并挂权限。

**Why this priority**: 商品是商城的根数据；没有后台可维护的目录，C 端无货可卖、交易链路无对象。

**Independent Test**: 建类目 → 建品牌 → 建 SPU（含 SKU）→ 列表/详情可见 → 上下架与限售生效 → 删除后不可见；全程按权限点校验。

**Acceptance Scenarios**:

1. **Given** 授权运营，**When** 调用类目/品牌/SPU/SKU 的 19 个管理端点，**Then** 全部正常返回（不再有"未实现"响应）。
2. **Given** 分类树请求，**When** 查询，**Then** 返回树形结构（父子嵌套）；创建子分类后父节点 children 含新节点。
3. **Given** 新货上架（含至少一个启用 SKU 与库存），**When** 置为下架，**Then** C 端列表不再返回该商品（既有规则：无启用 SKU 禁上架）。
4. **Given** SPU 配置限售区域，**When** 结算该商品至限售地区，**Then** 交易侧按既有规则拦截（005 已实现的限售校验不被本批回退）。
5. **Given** 不持 `product:spu:create` 的账号，**When** 调"新增商品"，**Then** 被拒（权限不足）；超管放行。

---

### User Story 2 - 库存后台（Priority: P1）

授权运营查询库存（按 SKU/关键词筛选，返回总量/锁定/可售推导）、查看预警（可售 ≤ 预警阈值）、执行库存调整（有符号增减，留痕：变更前后快照 + 操作人 + 备注）。

**Why this priority**: 库存是交易的前提（超卖防线）；后台无调整入口则运营无法补货/纠错。本故事是本批**唯一需要新写服务实现**的部分。

**Independent Test**: 造 SKU 库存（total=10/locked=2）→ 列表可售=8 → 调整 +5 后 total=15 且留痕一条 → 预警列表按阈值命中。

**Acceptance Scenarios**:

1. **Given** SKU 库存 total=10、locked=2，**When** 查询库存列表，**Then** 返回总量 10/锁定 2/可售 8（可售=总量−锁定）。
2. **Given** 库存总量 10、锁定 2、预警阈值 5，**When** 查询预警列表，**Then** 该 SKU 命中（可售 8 > 5 不命中则为其他阈值下命中——以可售 ≤ 阈值为判定）。
3. **Given** 某 SKU 库存 10，**When** 执行调整 +5（备注"补货"），**Then** 总量变 15，且变更流水记录"调整前 10/调整后 15/操作人/备注"。
4. **Given** 某 SKU 库存 10，**When** 执行调整 −3（备注"盘点亏损"），**Then** 总量变 7 并留痕。
5. **Given** 调整后总量将为负，**When** 提交，**Then** 拒绝（库存不可为负）——以既有规则为准。
6. **Given** 不持 `inventory:adjust` 的账号，**When** 调调整端点，**Then** 被拒；超管放行。

---

### User Story 3 - C 端商品浏览（Priority: P1）

游客/会员浏览商品：分类树、品牌列表、商品列表（含价格区间与可售标记）、商品详情（含 SKU 与详情字段）、商品搜索。服务实现已在 005 交付并有测试，本批接入契约端点（公开访问，无需登录）。

**Why this priority**: 浏览是 C 端入口；本批连线后商品数据才对前端可用。

**Independent Test**: 上架商品 → 游客（无凭证）调 5 个浏览端点 → 正常返回且仅见上架商品。

**Acceptance Scenarios**:

1. **Given** 无凭证游客，**When** 调分类树/品牌列表/商品列表/商品详情/搜索，**Then** 全部正常返回（公开）。
2. **Given** 已下架商品，**When** 游客浏览列表与搜索，**Then** 不出现；已下架商品详情返回商品不存在/失效错误。
3. **Given** 商品有多个 SKU 价格，**When** 浏览列表，**Then** 返回价格区间文本（既有口径不回退）。

---

### User Story 4 - 权限挂接（Priority: P2）

管理端 22 个端点按权限点校验：分类/品牌/SPU/SKU 各四码（read/create/update/delete）与 `inventory:read`/`inventory:adjust`；查询类端点仅要求登录；C 端 5 端点完全公开。

**Why this priority**: 延续批次 01~03 的 RBAC 机制；商品目录是核心数据，写操作必须受控。

**Independent Test**: 无权账号调各类写操作被拒；持 `product:spu:update` 但不持 `product:sku:create` 者改商品成功、建 SKU 被拒。

**Acceptance Scenarios**:

1. **Given** 不持 `product:spu:create` 的账号与超管，**When** 各自调"新增商品"，**Then** 前者被拒（权限不足）、超管成功。
2. **Given** 持 `product:spu:update` 但不持 `product:sku:create` 的账号，**When** 修改商品与新增 SKU，**Then** 前者成功、后者被拒（权限点各管一摊）。
3. **Given** 不持 `inventory:adjust` 的账号，**When** 调库存调整，**Then** 被拒。

### Edge Cases

- 既有实现的签名/字段与 api 契约不完全对齐时：以**契约（api Req/Res）为准**调整控制器映射，必要时最小改动 service 实现（记账）。
- 无 SKU 的 SPU 上架：拒绝（005 既有规则）。
- 分类删除有商品引用：拒绝（005 既有规则；错误码 30004）。
- 库存调整的 operator：按 inventory_log 表注释约定为 `admin:{id}`（非用户名）。
- 联动回退保护：本批不得回退交易链路（006）已验收的库存锁定/核销语义（不改其实现）。
- 本批零表结构变更、零新迁移（表齐备）。

## Requirements *(mandatory)*

### Functional Requirements

**后台商品目录（admin，19 端点，实现已有）**

- **FR-001**: 系统 MUST 提供分类管理：树查询、创建、修改、删除；删除被商品引用时拒绝。
- **FR-002**: 系统 MUST 提供品牌管理：列表（状态筛选+分页）、创建、修改、软删除。
- **FR-003**: 系统 MUST 提供 SPU 管理：列表（筛选+分页）、详情（含 SKU 与详情字段）、创建（返回编码）、修改、软删除、上下架切换、限售区域配置。
- **FR-004**: 系统 MUST 提供 SKU 管理：创建（挂 SPU）、修改、软删除、启停。
- **FR-005**: 上述 19 个端点 MUST 可调用（不再返回"未实现"），且行为与 005 已交付实现一致（不回退）。

**库存后台（admin，3 端点，实现新增）**

- **FR-006**: 系统 MUST 提供库存列表：按 SKU/关键词筛选 + 分页；返回总量/锁定/可售（可售=总量−锁定）与预警阈值。
- **FR-007**: 系统 MUST 提供预警列表：可售 ≤ 预警阈值的 SKU，分页。
- **FR-008**: 系统 MUST 提供库存调整：有符号增减；调整前后快照写入库存流水（含操作人、备注）；调整后总量 MUST NOT 为负。

**C 端浏览（shop，5 端点，实现已有）**

- **FR-009**: 系统 MUST 公开提供分类树、品牌列表、商品列表、商品详情、商品搜索；无需登录。
- **FR-010**: 游客浏览 MUST 只见上架商品（下架商品不可见）；搜索与列表同口径。
- **FR-011**: 商品列表 MUST 返回价格区间与可售标记（既有口径）。

**权限**

- **FR-012**: 管理端 22 端点 MUST 按权限点校验（分类/品牌/SPU/SKU 四码 + `inventory:read`/`inventory:adjust`）：持权放行、无权拒绝、超管直通；查询仅要求登录。

### Key Entities *(include if feature involves data)*

- **分类（product_category）**：父子层级、名称、层级、状态；树形展示。
- **品牌（product_brand）**：名称（存活期唯一）、状态、软删。
- **SPU（product_spu）**：编码（spu_no 唯一）、名称、类目/品牌、主图数组、规格定义、价格区间、上架状态、限售区域、软删。
- **SKU（product_sku）**：编码（sku_no 唯一）、所属 SPU、规格组合、价格（现价/划线价）、状态、软删。
- **库存（inventory）**：SKU 维度总量/锁定、预警阈值；可售=总量−锁定；变更留痕于库存流水。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 本批 27 个端点全部可调用（不再返回"未实现"），`make check-stub` 中 admin 渠道桩数由 91 降至 69、shop 渠道由 40 降至 35。
- **SC-002**: 库存三端点行为可验证：可售推导、预警阈值命中、调整留痕（行为级测试）。
- **SC-003**: 权限拦截可验证：无权写操作 100% 被拒、超管 100% 放行。
- **SC-004**: 回归零退化：005 既有商品测试与 006 交易链路测试全绿（本批不改其实现语义）。
- **SC-005**: `make test` 全绿，本批文件 lint 零问题。

## Assumptions

- **本批以"连线收口"为主**：19 个管理方法与 5 个浏览方法已在 005 交付（含测试），本批不改其业务逻辑；连线中发现签名/字段不匹配时按契约做最小调整并记账。
- 库存后台是**唯一新写实现**的部分（`IInventoryLogic` 的 List/Warnings/Adjust 三方法；该接口的交易侧方法由 006 以事务内条件更新实现，不在本批）。
- 库存调整的 operator 取登录后台账号（写入流水），实现期择一字段（用户名或 ID）并在契约注明。
- 预警阈值取自库存表既有列（warn_count），不改表结构。
- `IInventoryLogic` 等接口无实现者的问题（评审 I5 记录在案）不在本批处理。
- 批次 DoD：admin 91→69、shop 40→35（本批 27 端点清零）；全量测试绿；本批文件 lint 零问题。
