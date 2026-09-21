# Tasks: 会员中心（011-member-center）

**Input**: Design documents from `/specs/011-member-center/`

**Prerequisites**: plan.md、spec.md、research.md（D1~D6）、data-model.md、contracts/member-endpoint-mapping.md、quickstart.md

**Tests**: 宪法 IV + 分层契约 §五——**全部新写方法 TDD 红绿**（本批从零实现, 无既有测试兜底）。
测试落 `internal/service/user/*_impl_test.go`（复用该包 init 基座与 cleanState 惯例）。

**Organization**: 基建（迁移/DTO/接口）先行，随后按故事实现+连线（每故事先红后绿再接线）。

## Format: `[ID] [P?] [Story] Description`

- **[Story]**: US1 资料 / US2 地址 / US3 收藏足迹 / US4 消息偏好 / US5 积分 / US6 邀请
- 路径相对仓库根；后端命令在 `apps/server/` 下执行

## Phase 1: Setup（现状基线）

- [x] T001 基线确认：`go build ./...` 与 `go test ./...` 全绿（批次 04 收口态）；确认 user 桩数 34

## Phase 2: Foundational（阻塞前置，research D2/D5）

- [x] T002 [P] 新增迁移 `migrations/000036_user_notify_preference.{up,down}.sql`（会员×渠道唯一, D2）；
      `make migrate-fresh` 空库全量重放零失败（SC-004）；`make gen` 生成 dao/entity/do
- [x] T003 [P] `internal/model/dto_user.go` 新增 `InviteRecordItem`；`internal/service/user/distribution.go`
      的 `IDistributionLogic` 微扩 `InviteRecords`（D5）；PROGRESS 记账

## Phase 3: User Story 1 - 资料与登录记录 (Priority: P1) 🎯 MVP

**Goal**: 资料查询（脱敏）/修改/近 30 天登录记录（spec US1，FR-001~003）

**Independent Test**: 造会员 → 查资料（脱敏+等级+成长值）→ 改昵称 → 登录记录含最近一条

- [x] T004 [P] [US1] `internal/service/user/profile_impl_test.go`（红）：资料脱敏（明文不出参）/等级名按成长值匹配/
      改昵称后未传字段不变/登录记录 30 天窗口与分页倒序
- [x] T005 [US1] `internal/service/user/profile_impl.go`：`ProfileDetail`（脱敏+等级匹配）/`ProfileUpdate`/
      `LoginLogs`；同时实现 `GrowthAdd`/`LevelRecalc`（内部方法, D6）
- [x] T006 [US1] 连线 `controller/user/user_v1_profile_{detail,update}.go` + `user_v1_login_log_list.go`（桩清零 ×3）
- [x] T007 [US1] `go test ./internal/service/user/...` 绿

**Checkpoint**: 会员资料闭环。

---

## Phase 4: User Story 2 - 收货地址 (Priority: P1)

**Goal**: 地址 CRUD + 默认唯一 + 越权防护（spec US2，FR-004~006）

**Independent Test**: 建两条 → 设默认 → 唯一性成立 → 他人地址操作被拒

- [x] T008 [P] [US2] `address_impl_test.go`（红）：列表/新增/修改/删除；**设默认同事务清其他（无双默认）**；
      **他人地址改删设默认 → 10006**；`GetForOrder` 归属校验（内部方法）
- [x] T009 [US2] `address_impl.go`：5 端点方法 + `GetForOrder`
- [x] T010 [US2] 连线 5 端点（桩清零 ×5）
- [x] T011 [US2] `go test ./...` 绿

---

## Phase 5: User Story 3 - 收藏与足迹 (Priority: P1)

**Goal**: 收藏（实时价态/复活）+ 足迹（列表/清空）（spec US3，FR-007~008）

**Independent Test**: 收藏→列表价态→取消→再收藏复活→足迹清空

- [x] T012 [P] [US3] `favorite_impl_test.go`（红）：列表含实时价态（价格+可售）；**取消=软删**；
      **再次收藏=复活**（无重复行）；足迹列表倒序含次数；清空；`Record` UPSERT（内部）
- [x] T013 [US3] `favorite_impl.go`：收藏 3 + 足迹 2 端点方法 + `Record`/`CleanExpired`
- [x] T014 [US3] 连线 5 端点（桩清零 ×5）
- [x] T015 [US3] `go test ./...` 绿

---

## Phase 6: User Story 4 - 站内信与通知偏好 (Priority: P2)

**Goal**: 消息（未读计数/已读）+ 偏好（两渠道/默认全开）（spec US4，FR-009~011）

**Independent Test**: 2 未读→计数 2→单条已读→全读→计数 0；偏好默认全开→关闭→再查

- [x] T016 [P] [US4] `notify_impl_test.go`（红）：消息列表**未读计数（全量非当前页）**与已读态筛选；
      单条/全部已读；**偏好未设置=全开**；设置幂等（重复同值不重复行）；站内信不在偏好内
- [x] T017 [US4] `notify_impl.go`：消息 3 + 偏好 2 端点方法 + `Enqueue`/`DispatchTask`（内部, D6）
- [x] T018 [US4] 连线 5 端点（桩清零 ×5）
- [x] T019 [US4] `go test ./...` 绿

---

## Phase 7: User Story 5+6 - 积分与邀请 (Priority: P2/P3)

**Goal**: 积分账户/流水 + 邀请记录（spec US5/US6，FR-012~014）

**Independent Test**: 查账户与流水（类型筛选）；邀请仅本人

- [x] T020 [P] [US5] `point_impl_test.go`（红）：账户字段（余额可负如实）；流水类型筛选与分页；
      `Earn`/`Consume`/`Refund`/`ExpireDormant`（内部）；[US6] `invite_impl_test.go`（红）：仅本人记录+分页
- [x] T021 [US5] `point_impl.go`（账户/流水 + 内部四方法）+ `invite_impl.go`（邀请列表）
- [x] T022 [US5] 连线 3 端点：`user_v1_point_{account,log_list}.go` + `user_v1_invite_record_list.go`（桩清零 ×3）
- [x] T023 `go test ./...` 绿

**Checkpoint**: 21 端点全部去桩。

---

## Phase 8: Polish & 批次收尾（DoD 全项）

- [ ] T024 `make test` 全绿 + golangci-lint 本批文件零问题
- [ ] T025 `make check-stub` 对账：user 34→13（本批 21 清零）；`make migrate-fresh` 重放验证；
      按 quickstart 冒烟；更新 PROGRESS 批次 05 状态 ✅ 与完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001 → T002/T003 → US1 → US2 → US3 → US4 → US5+US6 → 收尾
- 各故事文件不重叠（profile/address/favorite/notify/point/invite），可按序高效推进
- 收尾依赖全部故事完成

## Implementation Strategy

1. 基线 → 迁移 + DTO/接口（生成物需 `make gen`）
2. 按故事：红 → 实现（含内部方法）→ 连线 → 绿 → 提交
3. 收尾：lint/桩数/迁移重放/冒烟 → PROGRESS ✅

## Notes

- 禁止手改 `internal/dao`、`internal/model/entity|do`（迁移后必须 `make gen` 重新生成）
- **越权防护是硬约束**：所有方法以会话 userId 收口, 他人资源按"不存在"处理（不泄露存在性）
- 手机号明文不得出参或落日志（个保法）
- 本批文件边界（PROGRESS 防偏离条款 3）：plan.md「Source Code」小节 + specs/011-member-center/ + specs/PROGRESS.md
