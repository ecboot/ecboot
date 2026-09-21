# Tasks: 物流公司与运营装修（009-logistics-ops）

**Input**: Design documents from `/specs/009-logistics-ops/`

**Prerequisites**: plan.md、spec.md、research.md（D1~D7）、data-model.md、contracts/operation-permissions.md、quickstart.md

**Tests**: 宪法 IV + 分层契约 §五 强制 TDD——先红后绿。测试落 `internal/service/shop/operation_impl_test.go`（复用同包 init 基座与清理惯例）。

**Organization**: 基建（DTO/接口/错误码）阻塞全部故事；US1 物流、US2 装修管理、US3 C 端展示可顺序推进。

## Format: `[ID] [P?] [Story] Description`

- **[Story]**: US1 物流字典 / US2 装修管理 / US3 C 端展示 / US4 权限
- 路径相对仓库根；后端命令在 `apps/server/` 下执行

## Phase 1: Setup（现状基线）

- [x] T001 基线确认：`go build ./...` 与 `go test ./...` 全绿（批次 02 收口态 + 废弃表已清）

**Checkpoint**: 基线绿。

---

## Phase 2: Foundational（阻塞前置，research D1/D5/D7）

- [x] T002 [P] `internal/model/dto_shop.go` 新增 7 个 DTO（域前缀命名）：`OperBannerItem/OperBannerInput/OperFloorItem/OperFloorInput/PublicBannerItem/PublicFloorItem/FloorProductSummary`（字段见 data-model §五与 api 契约）
- [x] T003 [P] `internal/errcode/errcode.go` 新增 `CodeLogisticsCodeTaken = 40012`（交易域段, D5）
- [x] T004 `internal/service/shop/misc.go`：`ILogisticsLogic` 微扩 `Detail`（D7）；新建 `IOperationLogic`（banner/floor 管理与 C 端共 10 方法, D1）；PROGRESS 变更记录记账

**Checkpoint**: DTO/接口/错误码就绪。

---

## Phase 3: User Story 1 - 物流公司字典 (Priority: P1) 🎯 MVP

**Goal**: 物流公司 CRUD + 编码唯一 + 停用保留（spec US1，FR-001~006）

**Independent Test**: 建公司 → 编码重复拒绝 → 停用仍在列表 → 软删后不可见

### Tests for US1（红先行）

- [x] T005 [P] [US1] `internal/service/shop/operation_impl_test.go`（物流部分）：创建成功/编码重复 40012/修改名称与状态（**编码不可改**——入参无该字段）/停用后列表仍可检出（status=0 筛选）/软删后列表与详情均 10006/详情不存在 10006/列表状态筛选与分页

### Implementation for US1

- [x] T006 [US1] `internal/service/shop/logistics_impl.go`：`LogisticsList`（status 筛选 + 分页, sort,id 序）、`LogisticsCreate`（编码唯一 + 必填校验）、`LogisticsUpdate`（名称/规则/排序/状态, 0 行→10006）、`LogisticsDelete`（软删, 0 行→10006）、`LogisticsDetail`
- [x] T007 [US1] 连线 `internal/controller/admin/admin_v1_admin_logistics_{list,create,update,delete,detail}.go`（桩清零 ×5；写操作挂 `logistics:company:manage`）
- [x] T008 [US1] `go test ./internal/service/shop/...` 绿

**Checkpoint**: 物流字典可用——发货选择的数据源就绪。

---

## Phase 4: User Story 2 - 运营装修管理 (Priority: P1)

**Goal**: banner/楼层 CRUD（spec US2，FR-007~009）

**Independent Test**: 建轮播（位置）→ 列表按位置筛选 → 建楼层（三型 + JSON 配置）→ 软删后不可见

### Tests for US2（红先行）

- [x] T009 [P] [US2] `operation_impl_test.go`（装修管理部分）：BannerCreate 成功/位置域外（service 兜底）/BannerList 位置筛选分页/修改含投放时段与启停/软删后管理端与 C 端均不可见；FloorCreate 成功（含 config JSON）/类型域外拒绝/修改 config 与启停/软删同理

### Implementation for US2

- [x] T010 [US2] `internal/service/shop/operation_impl.go`（管理部分）：`BannerList/BannerCreate/BannerUpdate/BannerDelete/FloorList/FloorCreate/FloorUpdate/FloorDelete`（位置与类型白名单校验, 0 行→10006, 软删）
- [x] T011 [US2] 连线 admin 8 端点：`admin_v1_admin_banner_{list,create,update,delete}.go`（挂 `operation:banner:manage`）、`admin_v1_admin_floor_{list,create,update,delete}.go`（挂 `operation:floor:manage`）（桩清零 ×8）
- [x] T012 [US2] `go test ./...` 绿（含批次 01/02 回归）

**Checkpoint**: 装修配置管理可用。

---

## Phase 5: User Story 3 - C 端装修展示 (Priority: P1)

**Goal**: 在投过滤 + 商品楼层摘要装配（spec US3，FR-010~013）

**Independent Test**: 三条轮播（在投/未到/过期）→ 仅见在投；商品楼层含失效 ID → 仅返回有效摘要

### Tests for US3（红先行）

- [x] T013 [P] [US3] `operation_impl_test.go`（C 端部分）：PublicBanners 时段三态（长期命中/未到排除/过期排除）+ 停用排除 + 排序；PublicFloors 启用楼层排序返回 + **商品楼层装配**（有效 SPU 出摘要 name/image/price_min）+ **失效剔除**（下架/不存在 ID 不报错）+ 停用排除 + 位置域外由 api 拦截的 service 兜底

### Implementation for US3

- [x] T014 [US3] `internal/service/shop/operation_impl.go`（C 端部分）：`PublicBanners`（在投 SQL 判定, D4）、`PublicFloors`（启用楼层 + 商品楼层装配 `spuIds` → SPU 摘要, 复用 `firstImage`, D2/D3, 失效剔除）
- [x] T015 [US3] 连线 `internal/controller/shop/shop_v1_banner_list.go`、`shop_v1_floor_list.go`（桩清零 ×2；公开, 白名单已有）
- [x] T016 [US3] `go test ./...` 绿

**Checkpoint**: 15 端点全部去桩。

---

## Phase 6: User Story 4 - 权限挂接 (Priority: P2)

- [x] T017 [P] [US4] `internal/middleware/permission_test.go` 补三码断言：无权账号对 `logistics:company:manage`/`operation:banner:manage`/`operation:floor:manage` 均拒 10005、超管放行；核对挂点与 contracts 一致（写操作挂、查询不挂、C 端公开）

**Checkpoint**: 权限语义闭合。

---

## Phase 7: Polish & 批次收尾（DoD 全项）

- [x] T018 `make test` 全绿 + golangci-lint 本批文件零问题（宪法 IV）
- [x] T019 `make check-stub` 对账：admin 104→91、shop 42→40（本批 15 清零）；按 quickstart 序列冒烟；更新 specs/PROGRESS.md 批次 03 状态 ✅ 与完成 commit（同 commit）并提交 `feat(009-logistics-ops): <收尾描述>`

## Dependencies & Execution Order

- T001 → T002/T003（并行）→ T004 → US1 → US2 → US3 → US4 → 收尾
- US1 与 US2/US3 不同文件（logistics_impl vs operation_impl）——可并行；US3 复用 US2 的实现文件故顺序执行
- 收尾依赖全部故事完成

## Implementation Strategy

1. 基线绿 → 基建（DTO/接口/错误码）
2. US1 物流字典（MVP：发货数据源）→ 独立验证
3. US2 装修管理 → US3 C 端展示（在投 + 装配）→ US4 权限断言 → 收尾数字对账
4. 每个任务或逻辑任务组完成即提交（`feat(009-logistics-ops): 中文描述`）

## Notes

- 禁止手改 `internal/dao`、`internal/model/entity|do`（生成物）
- 本批文件边界（PROGRESS 防偏离条款 3）：plan.md「Source Code」小节所列文件 + specs/009-logistics-ops/ + specs/PROGRESS.md；越界先记账再动
- 中文参数冒烟须 URL 编码（批次 02 教训）

## 完成记录（2026-09-21）

- 全量 `go test ./...` 绿（7 包）；本批文件 golangci-lint 0 issues（修 2 处 unused：shop 包 parseID 与 seedFloor）
- `make check-stub` 对账: admin 104→91、shop 42→40（本批 15 端点全清）；四渠道剩余 166 桩
- 冒烟 12/12: 物流 CRUD/40012/停用保留；轮播投放三态（长期 in/未到 out/过期 out）/位置域外 10001/停用 out；
  楼层装配（有效 SPU 出 name/image/price, 失效剔除, 下架剔除）/无权 10005×2/公开性
- 实现期发现: **既有 fixture `setupTradeFixture` 已失效**——归因更正（评审 M1）: `spu_no` 非空列自
  **000002** 起即存在（非 000034；000034 只加 price_min/price_max），该 fixture 从未写 spu_no。
  → 本批写自包含 fixture, 未修批次外既有 fixture（防偏离条款 3）
- 契约对齐: 轮播位置与楼层类型创建后不可改（api Update Req 无这两字段, spec 已修正）

## 评审修复轮（2026-09-21, No → 修）

- **C1（必修, 已修）**: `timeText`/`rbac_impl.rfc3339` 误用 gf 布局（`gtime.Format` 的 token 是 `Y-m-d H:i:s`，
  传 Go 布局产出字面串）→ 改标准库 `t.Time.Format(time.RFC3339)`；两处同源缺陷一并修
- **C2（必修, 已修——用户裁定方案 A）**: 时区口径不一致（Go 进程 +08 vs MySQL 会话 UTC）致强类型时间列
  读回偏移 8h、回显再提交每轮再漂 8h。修复：应用进程固定 UTC（`main.go` + 4 处测试基座
  `time.Local = time.UTC`），与库内 UTC 墙钟及会话时钟对齐；`TestTimeRoundTrip` 已由 Skip 转绿
- **I1（已修）**: service 层 status 白名单（`in.Status ∈ {0,1}`）三处 + 断言（域外值不落库）
- **I2（已修）**: 时段清空改 `g.DB().Transaction` 包裹（防半更新）+ 置空语句补 `deleted=0` +
  Count 错误显式传播（原先 `cerr == nil && cnt == 0` 会在 Count 出错时穿透并误报成功）
- **I3（已补）**: 断言补强——`TestClearTimeAndConfig`（Raw 置空/置 NULL 实证）、
  `TestSoftDeleteHidesPublic`（软删后 C 端不可见, FR-009）、`TestUpdateStatusGuard`（I1）、
  `TestTimeRoundTrip`（C2 灯）
- **I4（已修）**: `cleanupFloor` 改双路精确匹配（原名 + 改名后 `_2`）并清理实测残留 7 行；
  fixture 引用行前缀 `TF2-` → `TF-`（命中既有清理模式）并在清理时一并删
- **M1/M6（已记账）**: spu_no 归因更正；`api/admin/v1/logistics.go` 越界改动补记

### 修复轮 HTTP 复验（2026-09-21）

- C1: admin 出参 `startTime = 2026-12-31T16:00:00Z`（合法 RFC3339, 非 gf 布局字面串）✓
- C2: 同一瞬时校验通过; **回显原样再提交两次后仍为 `2026-12-31T16:00:00Z`（零漂移）** ✓
- I1: status=7 → `{"code":10001,"message":"状态须为1启用或0停用"}` ✓
