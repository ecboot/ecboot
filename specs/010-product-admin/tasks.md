# Tasks: 商品目录后台与 C 端浏览收口（010-product-admin）

**Input**: Design documents from `/specs/010-product-admin/`

**Prerequisites**: plan.md、spec.md、research.md（D1~D4）、data-model.md、contracts/product-endpoint-mapping.md、quickstart.md

**Tests**: 宪法 IV——**新写部分（库存三方法）TDD 红绿**；连线部分以"端到端可调用 + 权限 + 既有测试回归"验收（既有 005 测试承担语义回归，不重复造测试）。

**Organization**: 库存补齐（唯一新写）先行，随后两批连线（admin 19 + shop 5），最后权限断言与收尾。

## Format: `[ID] [P?] [Story] Description`

- **[Story]**: US1 后台目录连线 / US2 库存 / US3 C 端连线 / US4 权限
- 路径相对仓库根；后端命令在 `apps/server/` 下执行

## Phase 1: Setup（现状基线）

- [x] T001 基线确认：`go build ./...` 与 `go test ./...` 全绿（批次 03 收口态 + 时区统一）；确认 27 端点桩数 91/40

**Checkpoint**: 基线绿。

---

## Phase 2: User Story 2 - 库存后台 (Priority: P1) 🎯 MVP（唯一新写）

**Goal**: 库存列表（可售推导）/预警（阈值）/调整（留痕）（spec US2，FR-006~008）

**Independent Test**: 造 SKU 库存 total=10/locked=2 → 列表可售=8 → 调整 +5 → total=15 且流水一条 → 预警按阈值命中

### Tests for US2（红先行）

- [x] T002 [P] [US2] `internal/service/shop/inventory_impl_test.go`：List 可售推导（total-locked）与筛选分页；Warnings 阈值命中（available<=warn_count, 含等于）与不命中；Adjust 正负调整生效 + `inventory_log` 留痕（change_type=5/quantity=绝对值/total_after/locked_after/operator=admin:{id}/remark）；调整致负 → 30007；不存在 SKU → 30007

### Implementation for US2

- [x] T003 [US2] `internal/service/shop/inventory_impl.go`（新增）：`InventoryList`/`InventoryWarnings`/`InventoryAdjust`
      （JOIN product_sku 取名/编码；调整同事务写流水；条件更新防负, research D3）
- [x] T004 [US2] 连线 `internal/controller/admin/admin_v1_admin_inventory_{list,warn,adjust}.go`（桩清零 ×3；
      调整挂 `inventory:adjust`；查询仅登录）
- [x] T005 [US2] `go test ./internal/service/shop/...` 绿

**Checkpoint**: 库存后台可用（交易前提）。

---

## Phase 3: User Story 1 - 后台商品目录连线 (Priority: P1)

**Goal**: 19 个管理端点接入既有实现并挂权限（spec US1，FR-001~005）

**Independent Test**: 建类目/品牌/SPU/SKU → 列表详情可见 → 上下架与限售生效 → 权限拦截正确

### Implementation for US1（连线, 无新逻辑）

- [ ] T006 [US1] 连线分类 4 端点：`admin_v1_admin_category_{tree,create,update,delete}.go`
      （`shop.NewProductLogic()` 调用, research D1；写操作挂 `product:category:{create,update,delete}`）
- [ ] T007 [US1] 连线品牌 4 端点：`admin_v1_admin_brand_{list,create,update,delete}.go`（挂 `product:brand:*`）
- [ ] T008 [US1] 连线 SPU 7 端点：`admin_v1_admin_spu_{list,create,detail,update,delete,status,restrict}.go`
      （挂 `product:spu:{create,update,delete}`；status/restrict 挂 update）
- [ ] T009 [US1] 连线 SKU 4 端点：`admin_v1_admin_sku_{create,update,delete,status}.go`
      （挂 `product:sku:{create,update,delete}`；status 挂 update）
- [ ] T010 [US1] `go test ./...` 绿（含 005 商品测试回归 —— SC-004 零退化）

**Checkpoint**: 商品后台闭环。

---

## Phase 4: User Story 3 - C 端浏览连线 (Priority: P1)

**Goal**: 5 个浏览端点接入既有实现（公开）（spec US3，FR-009~011）

**Independent Test**: 无凭证调 5 端点正常；下架商品不出现在列表与搜索

### Implementation for US3（连线, 无新逻辑）

- [ ] T011 [US3] 连线 `shop_v1_category_tree.go`、`shop_v1_brand_list.go`、`shop_v1_product_list.go`、
      `shop_v1_product_detail.go`（viewerUserId 从 ctx 取, 未登录=0）、`shop_v1_product_search.go`
      （桩清零 ×5；公开——白名单已有）
- [ ] T012 [US3] `go test ./...` 绿 + 确认下架商品在列表/搜索不可见（复用既有 005 测试断言）

---

## Phase 5: User Story 4 - 权限断言 (Priority: P2)

- [ ] T013 [P] [US4] `internal/middleware/permission_test.go` 补商品域码断言：无权账号对
      `product:spu:create`/`product:sku:create`/`inventory:adjust` 等均拒 10005、超管放行；
      核对 22 个管理端点挂点与 contracts 一致（写操作挂、查询不挂、C 端公开）

---

## Phase 6: Polish & 批次收尾（DoD 全项）

- [ ] T014 `make test` 全绿 + golangci-lint 本批文件零问题
- [ ] T015 `make check-stub` 对账：admin 91→69、shop 40→35（本批 27 清零）；按 quickstart 冒烟；
      更新 specs/PROGRESS.md 批次 04 状态 ✅ 与完成 commit（同 commit）并提交 `feat(010-product-admin): <收尾>`

## Dependencies & Execution Order

- T001 → US2（库存新写）→ US1（admin 连线）→ US3（shop 连线）→ US4（权限断言）→ 收尾
- US1 与 US3 的 controller 文件不重叠（admin vs shop 包）——可交叉推进但建议顺序以控风险
- 收尾依赖全部故事完成

## Implementation Strategy

1. 基线绿 → 库存补齐（唯一新写, TDD）
2. admin 19 端点连线（按 contracts 映射表逐组推进）→ 005 回归验证
3. shop 5 端点连线 → 权限断言 → 收尾数字对账
4. 每个任务组完成即提交（`feat(010-product-admin): 中文描述`）

## Notes

- **不改既有实现语义**（005 商品 / 006 交易链路）——连线时若发现签名或字段不匹配, 按契约做最小适配并记账
- 禁止手改 `internal/dao`、`internal/model/entity|do`
- 本批文件边界（PROGRESS 防偏离条款 3）：plan.md「Source Code」小节 + specs/010-product-admin/ + specs/PROGRESS.md；越界先记账
