# Tasks: 后台账户与系统配置（007-admin-base）

**Input**: Design documents from `/specs/007-admin-base/`

**Prerequisites**: plan.md、spec.md、research.md（D1~D8 决策）、data-model.md、contracts/admin-permissions.md、quickstart.md

**Tests**: 宪法 IV + `docs/layer-contracts.md` §五 强制 TDD——每个故事先写失败测试（红）再实现（绿）。

**Organization**: 按用户故事分阶段； Foundational（D1 会话修复/权限挂接基建）阻塞全部故事。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 不同文件、无未完成依赖，可并行
- **[Story]**: US1 登录会话 / US2 账号管理 / US3 角色权限 / US4 系统配置 / US5 基础端点
- 路径均相对仓库根；后端命令在 `apps/server/` 下执行

## Phase 1: Setup（现状基线）

- [x] T001 基线确认：`go build ./...` 与 `go test ./...` 当前全绿，记录为回归基准（不改动任何文件）

**Checkpoint**: 基线绿——后续任何红都是本批引入的。

---

## Phase 2: Foundational（阻塞前置，research D1/D2/D3/D6）

**⚠️ CRITICAL**: 会话 audience 与权限基建未完成前，任何用户故事不得开始

- [x] T002 [P] 新增种子超管迁移 `apps/server/migrations/000035_admin_seed.up.sql`（INSERT is_super=1、status=1、bcrypt 哈希密码，实现期生成强初始密码）与 `000035_admin_seed.down.sql`（删除该行）；`make migrate-fresh` 空库重放零失败验证
- [x] T003 [P] `internal/errcode/errcode.go` 增补新错误码：权限不足/用户名已存在/角色编码已存在/存在引用禁删/禁止操作自身或超管/原密码错误/配置值类型不合法/生产环境禁用（语义见 contracts/admin-permissions.md）
- [x] T004 [P] `internal/model/dto_system.go` 新增 `AdminProfile`（Username/RealName/Roles []string）；`internal/service/system/rbac.go` 的 `IAdminAuthLogic` 微扩 `Profile(ctx, adminId int64) (*model.AdminProfile, error)`（research D6）
- [x] T005 会话 audience 维度（research D1）：`internal/library/security/session.go` 构造与 key 改为 `session:{aud}:{token}`/`session:refresh:{aud}:{token}`（aud: user/admin；`NewSessionManagerFromConfig` 适配）
- [x] T006 audience 调用点适配：`internal/middleware/auth.go`（按 `/admin/` 前缀判定渠道并校验对应会话；admin 会话额外校验 admin_user 存在且 status=1、deleted=0——账号被禁用/软删后会话立即不可用，spec FR-008）、`internal/service/user/auth.go` 会话创建点、`internal/controller/user/user_v1_logout.go`、`internal/controller/user/user_v1_token_refresh.go`
- [x] T007（随 T018 交付）新增 `internal/middleware/permission.go`：`RequirePerm(ctx, code) error`——取 CtxUserId 调 `system.HasPermission`，未持权返回权限不足错误码（research D2）
- [x] T008 回归适配：`internal/service/user/auth_flow_test.go` 适配 audience 签名；`go test ./...` 全绿（既有认证/交易回归通过）

**Checkpoint**: 会话双渠道隔离 + 权限挂接基建就绪；回归绿。

---

## Phase 3: User Story 1 - 后台登录与会话 (Priority: P1) 🎯 MVP

**Goal**: 种子超管可登录/登出/刷新/查资料/改密，登录全程审计（spec US1，FR-001~008）

**Independent Test**: 种子超管登录 → profile → 改密 → 旧密码被拒 → 登出双凭证失效，无需其他故事

### Tests for US1（红先行）

- [x] T009 [P] [US1] `internal/service/system/auth_impl_test.go`：登录成功返回双凭证+isSuper、登录审计成功行；密码错误拒绝+审计失败行(login_status=2)；账号不存在/禁用审计 login_status=3；验证码携带即校验（D7）；禁用账号登录被拒（FR-008）；改密后旧密码失效；Profile 返回角色编码列表

### Implementation for US1

- [x] T010 [US1] `internal/service/system/auth_impl.go`：`AdminLogin`（校验账号态/bcrypt/验证码/写 admin_login_log/更新 last_login_time/发 admin 会话）、`ChangePassword`（旧密码校验+bcrypt 新密码）、`Profile`（账号+角色联查）
- [x] T011 [US1] 连线 `internal/controller/admin/admin_v1_admin_login.go`、`admin_v1_admin_token_refresh.go`（复用 security.SessionManager admin 会话，user 渠道先例）、`admin_v1_admin_logout.go`、`admin_v1_admin_profile.go`、`admin_v1_admin_change_password.go`（桩清零 ×5）
- [x] T012 [US1] `go test ./internal/service/system/... ./internal/controller/...` 绿；US1 序列按 quickstart §手工验证 1/2/7 走查

**Checkpoint**: 后台可登录——MVP 可独立验证。

---

## Phase 4: User Story 2 - 后台账号管理 (Priority: P1)

**Goal**: 账号 CRUD/角色分配/列表筛选，含自我锁死与超管保护（spec US2，FR-009~012）

**Independent Test**: 创建→分配角色→列表可见→禁用后无法登录→软删后不可见

### Tests for US2（红先行）

- [x] T013 [P] [US2] `internal/service/system/rbac_impl_test.go`（账号部分）：创建成功/用户名冲突拒绝/修改姓名状态/禁用后 Login 被拒/软删后列表详情不可见且登录被拒/AssignRoles 全量替换幂等/删除自身与超管被拒（FR-012）/列表 status+keyword 筛选分页正确

### Implementation for US2

- [x] T014 [US2] `internal/service/system/rbac_impl.go`（账号部分）：`AdminUserList`（status/keyword/分页, 关联角色编码+last_login_time）、`AdminUserCreate`（用户名唯一+bcrypt）、`AdminUserUpdate`（real_name/status）、`AdminUserDelete`（软删；禁删自身/超管）、`AssignRoles`（事务删旧插新）
- [x] T015 [US2] 连线 `internal/controller/admin/admin_v1_admin_{user_list,user_create,user_update,user_delete,user_detail,user_assign_roles}.go`（桩清零 ×6；权限挂点随 US3 的 RequirePerm 就绪后补挂——编译依赖顺序, 见 T019）
- [x] T016 [US2] `go test ./...` 绿（含 US1 回归）

**Checkpoint**: US1+US2 独立可用。

---

## Phase 5: User Story 3 - 角色与权限管理 (Priority: P2)

**Goal**: 角色 CRUD/权限树/全量替换授权/权限拦截生效（spec US3，FR-013~018）

**Independent Test**: 建角色→分配权限→挂到新账号→持权放行/无权拒绝/超管直通

### Tests for US3（红先行）

- [x] T017 [P] [US3] `internal/service/system/rbac_impl_test.go`（角色部分）：角色创建/编码冲突拒绝/修改/有账号引用禁删（FR-014）/无引用可软删/RoleDetailView 含 permissionIds/PermissionTree 与 000032 种子同源同层级（children 嵌套）/AssignPermissions 全量替换幂等/HasPermission：is_super 直通、持权 true、无权 false、权限 status=0 不命中（FR-017/018）

### Implementation for US3

- [x] T018 [US3] `internal/service/system/rbac_impl.go`（角色部分）：`RoleList/RoleCreate/RoleUpdate/RoleDelete/RoleDetailView/PermissionTree/AssignPermissions/HasPermission`（判定数据流见 data-model §三）
- [x] T019 [US3] 连线 `internal/controller/admin/admin_v1_admin_{role_list,role_create,role_update,role_delete,role_detail,role_assign_perm,permission_tree}.go`（桩清零 ×7；写操作挂对应权限点 `system:role:manage`/`system:role:assign`）
- [x] T020 [US3] 权限拦截行为验证（spec SC-003）：非超管无权账号调 `POST /admin/roles` 被拒、超管放行的集成断言（service 层 RequirePerm 语义级测试）
- [x] T021 [US3] `go test ./...` 绿（含 US1/US2 回归）

**Checkpoint**: RBAC 从数据管理变为真实访问控制。

---

## Phase 6: User Story 4 - 系统配置 (Priority: P2)

**Goal**: 配置列表/修改/类型校验，覆盖层停用语义（spec US4，FR-019~021）

**Independent Test**: 改值→列表反映→停用→消费方回退默认（research D5：读取组件不在本批）

### Tests for US4（红先行）

- [x] T022 [P] [US4] `internal/service/system/config_impl_test.go`：List 返回 000030 种子全项（含说明/默认值语义）/Update 改值与状态/整数型配置传非数字拒绝（valueType 校验 FR-020）/停用后 List 返回 status=0（覆盖层语义 FR-021）

### Implementation for US4

- [x] T023 [US4] `internal/service/system/config_impl.go`：`List`（全量, 按 sort/code 序）、`Update`（value 按类型校验 + status 切换）
- [x] T024 [US4] 连线 `internal/controller/admin/admin_v1_admin_{config_list,config_update}.go`（桩清零 ×2；Update 挂 `system:config:update`）
- [x] T025 [US4] `go test ./...` 绿

**Checkpoint**: US1~US4 全部独立可用。

---

## Phase 7: User Story 5 - 基础端点 (Priority: P3)

**Goal**: ping 探活 + mock 短信调试桩（spec US5，FR-022/023）

**Independent Test**: ping 无凭证 200；mock_latest_sms 非生产返回、生产拒绝

### Implementation for US5

- [x] T026 [P] [US5] 连线 `internal/controller/common/common_v1_ping.go`：直接返回存活 Res（无业务依赖；桩清零 ×1）
- [x] T027 [P] [US5] `internal/service/system/devtools_impl.go` 导出 `MockLatestSms(ctx, phone)`：非生产环境读 Redis `mock:sms:{phone}` 返回（`library/sms` 落点约定），生产拒绝（research D4）；连线 `internal/controller/common/common_v1_mock_latest_sms.go`（桩清零 ×1）；环境判定与生产拒绝的单元测试

**Checkpoint**: 22 端点全部去桩。

---

## Phase 8: Polish & 批次收尾（DoD 全项）

- [ ] T028 `make test` 全绿（含认证/交易既有回归）+ `make lint` 通过（宪法 IV）
- [ ] T029 `make check-stub` 对账：admin 129→107、common 5→3（本批 22 清零）；结果记入 PROGRESS.md 批次 01 状态 ✅ 与完成 commit（同 commit 更新）
- [ ] T030 按 `specs/007-admin-base/quickstart.md` 手工序列走查一遍；更新 specs/PROGRESS.md 批次 01 行（状态/commit）并提交 `feat(007-admin-base): <收尾描述>`

## Dependencies & Execution Order

- **Phase 1 → 2**：基线绿后动基建
- **Phase 2 阻塞全部故事**（T005~T008 会话/权限基建；T002 种子是 US1 登录测试前置）
- **US1 → US2 → US3**：US2 复用 US1 的 Login 断言禁用语义；US3 的 HasPermission 支撑 US2 已挂的 RequirePerm 端到端生效
- **US4、US5 与 US2/US3 无依赖**，可在基建完成后并行
- **Phase 8 收尾**依赖全部故事完成

### Parallel Opportunities

- Phase 2 内 T002/T003/T004 互不依赖可并行
- US4、US5 全程可与其他故事并行（不同文件）
- 各故事内测试任务（红）可先于实现批量编写

## Implementation Strategy

1. Phase 1+2 完成后 STOP：确认回归绿（安全修复波及既有渠道，必须先站稳）
2. US1 交付即 MVP（后台可进入）→ 独立验证
3. US2→US3 顺序交付（权限体系闭环）→ 每故事 Checkpoint 独立验证
4. US4/US5 穿插并行补齐 → Phase 8 DoD 数字收口
5. 每个任务或逻辑任务组完成即提交（中文 Conventional Commits，前缀 `feat(007-admin-base):`）

## Notes

- 禁止手改 `internal/dao`、`internal/model/entity|do`（生成物）
- 全部分层行为遵守 `docs/layer-contracts.md`（controller 禁触 dao；service 用 do/entity 强类型）
- 测试基座：005/006 模式（确定性配置注入 + 数据自建清理）
- 本批文件边界（PROGRESS 防偏离条款 3）：plan.md「Source Code」小节所列文件 + specs/007-admin-base/ + specs/PROGRESS.md；越界先记账再动
