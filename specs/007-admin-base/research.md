# Research: 007-admin-base 后台账户与系统配置

> Phase 0 产出。每个决策含：结论 / 理由 / 已评估替代方案。全部基于代码现状勘察（commit 9c8a74e 时点）。

## D1 会话增加 audience（渠道）维度 —— 安全修复

**结论**：`security.SessionManager` 的 Redis key 从 `session:{token}` / `session:refresh:{token}` 改为带渠道维度
`session:{aud}:{token}` / `session:refresh:{aud}:{token}`（aud ∈ `user` / `admin`）；构造函数增加 audience 参数；
`middleware.Auth` 按路径前缀判定渠道（`/admin/` → admin，其余会员/商城路径 → user）并校验对应 audience 的会话。

**理由**：现状 key 仅含 userId，admin 与 user 两表主键空间独立——user id=7 的合法 token 调 admin 接口时，
若 admin_user 恰有 id=7 记录即可冒充该管理员（跨渠道越权）。这是 FR-017 权限体系的前置安全缺陷，必须本批修复。

**替代方案**：Validate 后回查 admin_user 表确认身份——无法根治（token 只存 userId，无法判断它属于哪张表），
仅缩小碰撞面；放弃。

**波及面**（编译器兜底 + 既有测试回归）：`middleware/auth.go`、user 渠道 `logout`/`token_refresh` 控制器、
`service/user/auth.go` 会话创建点、`auth_flow_test.go`、`trade_test` 无涉。

## D2 权限点挂接 = controller 显式 helper（非中间件路由映射）

**结论**：新增 `middleware.RequirePerm(ctx, code) error`——从 ctx 取管理员 ID，调 `system.HasPermission`；
受权限点保护的 controller 方法首行显式调用。权限点编码以 api 层注释与 000032 种子为准绳。

**理由**：GoFrame 中间件内拿到的是真实路径（`/roles/9`），做 pattern 归一化（`{id}` 泛化）脆弱难测；
显式调用点即文档、TDD 友好、评审可见（漏挂靠评审清单而非运行时兜底）。
`HasPermission` 数据流：admin_user.is_super 直通 → admin_user_role → admin_role_permission → admin_permission.code。

**替代方案**：中间件内维护 `METHOD + pattern → 权限码` 映射表——路径归一化复杂、与 gf 路由 pattern 双份维护；放弃。

## D3 种子超管账号 = 新增迁移 000035_admin_seed（spec 假设修订）

**结论**：新增 `migrations/000035_admin_seed.{up,down}.sql`：up 插入 1 个 `is_super=1` 启用账号（bcrypt 哈希，
实现期生成，初始密码满足复杂度并在部署文档标注"首次登录必改"）；down 删除该行。**无表结构变更，纯数据种子**（000032 先例）。

**理由**：勘察发现 000010/000032 均无种子管理员——后台没有可登录账号，整个 admin 渠道无法进入，本批"后台可用"目标不成立。
spec Assumptions 原写"无新增迁移"，按 PROGRESS 防偏离条款 4 在 spec 变更记录记账后修订。

**替代方案**：测试自建 + 生产引导 SQL 放 docs 由运维手工执行——部署面易遗漏，种子走迁移可被 `migrate-fresh` 全量重放验证；放弃。

## D4 MockLatestSms = 运行时环境判定

**结论**：实现 MockLatestSms 控制器：读 Redis `mock:sms:{phone}`（`library/sms` 的 MockSender 落点）返回最近验证码；
环境为生产（配置 `system.env` 或等效为 prod）时返回拒绝错误，不泄露任何内容。

**理由**：api 注释期望"仅 mock 模式注册路由"，但路由是 `Bind(common.NewV1())` 全量绑定（gf gen ctrl 聚合接口），
拆单绑定要动生成物治理面。运行时判定同样满足 spec FR-023（生产拒绝且不泄露），代价最小。

**替代方案**：路由拆绑 mock 端点——动 `gf gen ctrl` 聚合接口与治理约定，收益仅是生产少一条 404 路由；放弃。

## D5 system_config 读取组件推迟到首个消费批次

**结论**：本批只交付 `IConfigLogic.List/Update` 管理面；不新建 library/config 读取组件。

**理由**：勘察确认本批 22 端点无业务配置消费方（`NewSessionManagerFromConfig` 已自带 SQL 直读回退先例）。
FR-021 的"覆盖层回退"行为由：① Update 正确写 status=0/1；② 既有 FromConfig 回退测试佐证。组件待 008（门店自提开关等）落地。

**替代方案**：预防性建读取组件——无消费方的抽象违反宪法 V（YAGNI）；放弃。

## D6 Profile 查询微扩 IAdminAuthLogic

**结论**：`IAdminAuthLogic` 增加 `Profile(ctx, adminId int64) (*model.AdminProfile, error)`（返回用户名/姓名/角色编码列表）。
spec Assumptions 已预告，DTO 补入 `model/dto_system.go`。

**理由**：AdminProfile 端点需要账号 + 角色联查，controller 禁触 dao（分层契约 §四），必须经 service。

## D7 登录验证码：携带即强校验

**结论**：`CaptchaKey/CaptchaCode` 非空时必校验（复用 `library/captcha` 校验能力），错误拒绝；为空跳过。
强制携带策略由后续风控/配置决定。

## D8 登录失败不锁定，仅审计留痕

**结论**：Login 失败写 `admin_login_log`（login_status：2 密码错误 / 3 账号禁用或不存在），不锁定账号。
锁定机制属风控域（018 批之后），表结构亦无锁定字段。

## 勘察结论（非决策，供 plan/tasks 引用）

- **DTO 就位**：`model/dto_system.go` 已定义 RoleItem/RoleInput/RoleDetailView/PermissionNode/AdminUserItem/
  AdminUserInput/AdminUserUpdateInput/AdminLoginResult/ConfigItem——本批零新增（仅 D6 的 AdminProfile）。
- **service 形态**：包级导出函数（`user.SmsLogin` 先例），接口定义作契约文档；controller 直接调 `system.Xxx(...)`。
- **错误码**：`internal/errcode` 统一出口（10003 未登录、20001 验证码错误等既有码复用；权限不足等新码按既有风格增补）。
- **admin 五表**：均有 `deleted` 软删列与 `status` 列；`admin_user_role`/`admin_role_permission` 为纯关联表。
- **白名单现状**：`/admin/login`、`/admin/token/refresh` 已在 Auth 白名单；`/common/` 全公开（ping 直通成立）。
- **测试基座**：005/006 模式（确定性配置注入 + 数据自建清理），`service/system` 包内新建 `*_test.go`。
