# Tasks: 门店域（008-store）

**Input**: Design documents from `/specs/008-store/`

**Prerequisites**: plan.md、spec.md、research.md（D1~D5）、data-model.md、contracts/store-permissions.md、quickstart.md

**Tests**: 宪法 IV + 分层契约 §五 强制 TDD——先红后绿。测试落 `internal/service/shop/store_impl_test.go`（复用该包既有 init 基座）。

**Organization**: 按用户故事分阶段；Foundational（idgen/AdminDetail 微扩）阻塞故事。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 不同文件、无未完成依赖，可并行
- **[Story]**: US1 游客检索 / US2 后台管理 / US3 权限语义
- 路径相对仓库根；后端命令在 `apps/server/` 下执行

## Phase 1: Setup（现状基线）

- [x] T001 基线确认：`go build ./...` 与 `go test ./...` 当前全绿（批次 01 收口态），记录为回归基准

**Checkpoint**: 基线绿。

---

## Phase 2: Foundational（阻塞前置，research D2/D3）

**⚠️ CRITICAL**: idgen 与接口微扩未完成前，用户故事不得开始

- [x] T002 [P] 新增 `internal/library/idgen/idgen.go`：sonyflake 单例 + `NextID() (uint64, error)`（先写唯一性测试 `idgen_test.go`：连续两次 NextID 不相等——红→绿）
- [x] T003 [P] `internal/service/shop/misc.go` 的 `IStoreLogic` 微扩 `AdminDetail(ctx, storeId int64) (*model.StoreItem, error)`（research D3, 同批次 01 D6 模式）；变更记录记账

**Checkpoint**: 编码生成与接口面就绪。

---

## Phase 3: User Story 1 - 游客门店检索 (Priority: P1) 🎯 MVP

**Goal**: 游客按区县或附近检索营业门店、看详情（spec US1，FR-001~006）

**Independent Test**: 造 3 店（2 营业 1 歇业，不同区县坐标）→ 区县筛选仅营业 → 附近检索距离升序且歇业/无坐标不出现 → 详情完整

### Tests for US1（红先行）

- [x] T004 [P] [US1] `internal/service/shop/store_impl_test.go`：PublicList 区县筛选仅营业/附近检索距离升序+distanceM/半径过滤（默认 10km、上限 100 回退）/无坐标店排除/区县+经纬度同传附近优先/PublicDetail 歇业店可见且不存在返回 10006；种子门店自建清理

### Implementation for US1

- [x] T005 [US1] `internal/service/shop/store_impl.go`：`PublicList`（包围盒预筛 + Haversine 计算列 + HAVING + 距离排序，research D1；区县模式 sort,id 排序）、`PublicDetail`（status/deleted 过滤）
- [x] T006 [US1] 连线 `internal/controller/common/common_v1_store_list.go`、`common_v1_store_detail.go`（桩清零 ×2；ID string↔int64 转换走 parseID 同款语义）
- [x] T007 [US1] `go test ./internal/service/shop/... ./internal/controller/...` 绿

**Checkpoint**: 游客可找店——MVP 可独立验证。

---

## Phase 4: User Story 2 - 后台门店管理 (Priority: P1)

**Goal**: 门店 CRUD + 编码生成 + 列表筛选（spec US2，FR-007~012）

**Independent Test**: 创建（返回 ST 编码）→ 列表可见 → 歇业后游客列表消失 → 软删后全端点不可见

### Tests for US2（红先行）

- [x] T008 [P] [US2] `internal/service/shop/store_impl_test.go`（管理部分）：AdminCreate 生成 ST 前缀唯一编码/区划码非 6 位数字拒绝（10001）/AdminUpdate 含歇业切换与目标不存在 10006/AdminDelete 软删后各端点不可见/AdminList status+keyword（名称或编码）筛选分页/AdminDetail 返回编码与坐标

### Implementation for US2

- [x] T009 [US2] `internal/service/shop/store_impl.go`（管理部分）：`AdminCreate`（校验 + idgen 编码 + 冲突重试 ≤3）、`AdminUpdate`（0 行=10006）、`AdminDelete`（软删）、`AdminList`（status/keyword/分页, sort,id 排序）、`AdminDetail`
- [x] T010 [US2] 连线 `internal/controller/admin/admin_v1_admin_store_{list,create,update,delete,detail}.go`（桩清零 ×5；写操作挂 `middleware.RequirePerm(ctx, "store:manage:create/update/delete")`，见 contracts）
- [x] T011 [US2] `go test ./...` 绿（含批次 01 回归）

**Checkpoint**: US1+US2 独立可用。

---

## Phase 5: User Story 3 - 权限语义 (Priority: P2)

**Goal**: store 权限点拦截生效 + 歇业语义数据自洽（spec US3，FR-013/014）

**Independent Test**: 无权账号写操作被拒、超管放行、持 update 不持 delete 者改得删不得

### Tests for US3（红先行）

- [x] T012 [P] [US3] `internal/middleware/permission_test.go` 补 store 码断言：无权账号对 `store:manage:create/update/delete` 均拒（10005）、超管放行（RequirePerm 语义级, 沿用批次 01 测试形态）

### Implementation for US3

- [x] T013 [US3] 确认 T010 挂点与 contracts 一致（create/update/delete 三处、查询不挂）；歇业语义由 T004/T008 断言覆盖核对（游客列表仅营业 + 详情 status 如实返回）
- [x] T014 [US3] `go test ./...` 绿

**Checkpoint**: 7 端点全部去桩且权限语义闭合。

---

## Phase 6: Polish & 批次收尾（DoD 全项）

- [x] T015 `make test` 全绿（含既有回归）+ golangci-lint 本批文件零问题（宪法 IV）
- [x] T016 `make check-stub` 对账：admin 109→104、common 3→1；按 quickstart 序列冒烟；更新 specs/PROGRESS.md 批次 02 状态 ✅ 与完成 commit（同 commit）并提交 `feat(008-store): <收尾描述>`

## Dependencies & Execution Order

- T001 → T002/T003（并行）→ US1 → US2 → US3 → 收尾
- US1 与 US2 同文件（store_impl.go）——顺序执行避免冲突；US3 测试可与 US2 并行（不同文件）
- 收尾依赖全部故事完成

### Parallel Opportunities

- T002/T003 互不依赖可并行
- T012（middleware 测试）与 US2 实现可并行

## Implementation Strategy

1. 基线绿 → 基建（idgen/接口微扩）
2. US1 交付即 MVP（游客可检索）→ 独立验证
3. US2 管理面闭环 → US3 权限断言收口 → Phase 6 数字收口
4. 每个任务或逻辑任务组完成即提交（`feat(008-store): 中文描述`）

## Notes

- 禁止手改 `internal/dao`、`internal/model/entity|do`（生成物）
- 全部分层行为遵守 `docs/layer-contracts.md`
- 本批文件边界（PROGRESS 防偏离条款 3）：plan.md「Source Code」小节所列文件 + specs/008-store/ + specs/PROGRESS.md；越界先记账再动

## 完成记录（2026-09-21）

- 全量 `go test ./...` 绿（7 包）；本批文件 golangci-lint 0 issues（修复 QF1008；批次外 15 个既有问题留 chore 批）
- `make check-stub` 对账: admin 109→104、common 3→1（本批 7 端点全清）；四渠道剩余 181 桩
- HTTP 冒烟 10/10: 超管建店(ST sonyflake 编码)/后台详情/游客区县检索(无认证)/附近检索(distanceM=147 与坐标吻合
  )/游客详情/无权建店 10005/后台list仅需登录/置歇业/歇业后游客列表清零
- 环境事项: 测试库补建缺失的 store 表（000031 未真正应用, 历史迁移改写致账本失真；非破坏性补建,
  merchant/seller/shop 三个废弃表残留未删——DROP 已被安全策略拦截, 待用户裁定）
- 记录在案: 非法编码入参返回 52 而非 10001（横切层既有行为, 跨批次）

## 评审修复轮（2026-09-21, With fixes → 已闭合）

- **C1（必修）**: AdminUpdate 状态白名单（仅 1营业/2歇业）+ api 层 `v:"required|in:1,2"`——防状态机外值落库致门店静默消失
- **I1**: 全量覆盖必填护栏（service 兜底 10001 + api required）——防"部分提交静默清空档案"（D6 取舍的前提落地）
- **I2**: 坐标清零改 `gdb.Raw("NULL")` 显式置空——实证 gdb 对 do 自动 OmitNilData 丢弃 nil（Fields 白名单亦救不回）
- **I3**: 附近检索分页双向钳制（`start<0` 溢出）——公开端点抗畸形入参
- **I6**: idgen 机器码 hostname FNV 回退——防"无可用私网 IP 时一次失败被 sync.Once 永久缓存"
- **M2/M4/M6**: 距离排序补 `id ASC` 兜底 / 经纬度域校验（越界返 10001）/ `go mod tidy` 摘 indirect
- **补测试**: 6 项新断言（守卫/坐标清零/分页边界/域校验/单坐标回退/并发唯一性）
- **待裁定（跨批次形态, 不在本批）**:
  - **I4**: response.mapError 对任何 DB 错误暴露 code 52 且绕过 g_log 与追踪号（10002 语义失效）——影响批次 01/02 全部 DB 错误路径, 单开 chore 批
  - **I5**: `IStoreLogic` 等接口无实现者（service 实现为包级函数, 与批次 01 system 域同形态）——「删接口（宪法 V 简单优先）」vs「补 struct 实现」二选一, 需与 `docs/layer-contracts.md` §七 同步决策
