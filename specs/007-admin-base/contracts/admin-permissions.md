# Contract: 007-admin-base 端点 × 权限点映射

> 22 端点全景（spec SC-001 对账基准）。权限列为 controller 显式 `RequirePerm` 挂点（research D2）；
> 「公开」= Auth 白名单直通；「登录」= 有效 admin 会话即可。错误码沿用 `internal/errcode` 既有风格。

## admin 渠道（20）

| 端点 | 方法+路径 | 鉴权 | 权限点 |
|---|---|---|---|
| AdminLogin | POST /admin/login | 公开 | — |
| AdminTokenRefresh | POST /admin/token/refresh | 公开 | — |
| AdminLogout | POST /admin/logout | 登录 | — |
| AdminProfile | GET /admin/profile | 登录 | — |
| AdminChangePassword | PUT /admin/profile/password | 登录 | — |
| AdminUserList | GET /admin/admin-users | 登录 | — |
| AdminUserCreate | POST /admin/admin-users | 登录 | system:admin:manage |
| AdminUserUpdate | PUT /admin/admin-users/{id} | 登录 | system:admin:manage |
| AdminUserDelete | DELETE /admin/admin-users/{id} | 登录 | system:admin:manage |
| AdminUserDetail | GET /admin/admin-users/{id} | 登录 | — |
| AdminUserAssignRoles | PUT /admin/admin-users/{id}/roles | 登录 | system:admin:assign |
| AdminRoleList | GET /admin/roles | 登录 | — |
| AdminRoleCreate | POST /admin/roles | 登录 | system:role:manage |
| AdminRoleUpdate | PUT /admin/roles/{id} | 登录 | system:role:manage |
| AdminRoleDelete | DELETE /admin/roles/{id} | 登录 | system:role:manage |
| AdminRoleDetail | GET /admin/roles/{id} | 登录 | — |
| AdminRoleAssignPerm | PUT /admin/roles/{id}/permissions | 登录 | system:role:assign |
| AdminPermissionTree | GET /admin/permissions | 登录 | — |
| AdminConfigList | GET /admin/configs | 登录 | — |
| AdminConfigUpdate | PUT /admin/configs/{code} | 登录 | system:config:update |

> 权限点编码与 api 层注释一致（`api/admin/v1/{auth,rbac,config}.go`）。
> 查询类端点仅要求登录（可读即可见）；写操作挂权限点——与 004-api-surface 注释约定对齐。

## common 渠道（2）

| 端点 | 方法+路径 | 鉴权 | 说明 |
|---|---|---|---|
| Ping | GET /common/ping | 公开（/common/ 全公开） | 存活探针 |
| MockLatestSms | GET /common/mock-latest-sms | 公开 | 生产环境拒绝（research D4） |

## 本批新增错误码（沿 errcode 既有风格）

| 码 | 语义 | 触发 |
|---|---|---|
| 10005（复用） | 权限不足 | RequirePerm 未持权且非超管，不新增登记 |
| 10006 | 资源不存在 | 通用码（新增通用段）：配置不存在/模拟短信无记录 |
| 80001（既有） | 后台凭证错误 | AdminLogin 用户名不存在/密码错误 |
| 80009 | 后台账号已禁用 | AdminLogin 命中 status=2 或 deleted=1 账号 |
| 80002（既有） | 原密码错误 | AdminChangePassword |
| 80003 | 用户名已存在 | AdminUserCreate 冲突 |
| 80004 | 角色编码已存在 | AdminRoleCreate 冲突 |
| 80005 | 存在引用禁删 | AdminRoleDelete 被 admin_user_role 引用 |
| 80006 | 禁止操作自身/超管 | AdminUserDelete(id=自身 或 is_super) |
| 80007 | 配置值类型不合法 | AdminConfigUpdate valueType 校验失败 |
| 80008 | 生产环境禁用 | MockLatestSms 于 prod（调试端点） |
