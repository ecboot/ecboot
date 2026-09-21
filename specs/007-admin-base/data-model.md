# Data Model: 007-admin-base

> 零表结构新增（D3 种子迁移为纯数据）。本文档摘录本批消费的既有结构与新增语义，
> 权威定义以迁移文件与 gf 生成物为准。

## 一、实体与消费字段

### admin_user（后台账号）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| id | 主键 | 会话 userId、路径参数 |
| username | 登录名 | **唯一**（创建冲突拒绝；查询按精确匹配） |
| password | 密码 | bcrypt 哈希，任何环节不出明文 |
| real_name | 姓名 | 创建可填、修改可改 |
| is_super | 超管标记 | 1=跳过权限点校验；**禁删** |
| status | 1正常 2禁用 | 禁用不可登录 |
| last_login_time | 最后登录 | 登录成功时更新 |
| deleted | 软删 | 1=不可登录、列表不可见 |

### admin_role（角色）
code **唯一**、status 1启用/0停用、deleted 软删；**被 admin_user_role 引用时禁删**（FR-014）。

### admin_permission（权限）
code=`模块:资源:动作`、type 1菜单/2按钮/3接口、parent_id 层级（0=根）、sort；树与 000032 种子同源。

### admin_user_role / admin_role_permission（关联表）
纯 M:N；分配 = **全量替换**（同事务删旧插新，重复执行幂等）。

### admin_login_log（登录审计，只追加）
login_status：1 成功 / 2 密码错误 / 3 账号禁用或不存在；记录 IP、User-Agent。

### system_config（系统配置，000030 已种子）
value_type：1整数/2小数/3字符串/4布尔/5JSON（Update 按类型校验值）；status 1启用/0停用
（停用=业务读取方回退代码默认，见 research D5）。

## 二、会话对象（Redis，本批改形）

```
session:{aud}:{token}           = userId    （访问凭证, 滑动续期, TTL=配置天数默认7）
session:refresh:{aud}:{token}   = userId    （刷新凭证, TTL=访问4倍, GETDEL 一次性）
aud ∈ { "user", "admin" }                    —— 新增维度（research D1）
```

判定规则：`/admin/` 前缀路径 → admin 会话；`/user/`、`/shop/` → user 会话；互不通用。

## 三、权限判定数据流（FR-017/018）

```
RequirePerm(ctx, code)
  → adminId = CtxUserIdFrom(ctx)
  → system.HasPermission(ctx, adminId, code)
      ├─ admin_user: deleted=0 且 status=1 且 is_super=1 → true（直通）
      └─ JOIN admin_user_role → admin_role_permission → admin_permission(code, status=1, deleted=0) 命中 → true
  → 否则 false → 统一权限不足错误码
```

## 四、状态迁移（本批范围）

- admin_user：创建(1) → 禁用(2) → 启用(1)；任意态 → 软删(1, deleted=1)（终态，不可登录/不可见）。
- admin_role：启用(1) ↔ 停用(0)；deleted=1 终态（有引用禁删）。
- system_config：启用(1) ↔ 停用(0)（覆盖层开关，无删除语义）。

## 五、新增 DTO（仅 1 个）

`model.AdminProfile`：Username / RealName / Roles []string（research D6）——补入 `dto_system.go`。
其余（RoleItem/RoleInput/RoleDetailView/PermissionNode/AdminUserItem/AdminUserInput/AdminUserUpdateInput/
AdminLoginResult/ConfigItem）全部复用既有定义。
