# Research: 全角色 API 接口面（Phase 0）

2026-09-20。6 项契约约定决策（全部有用户既有代码或 003 契约为锚点）。

## D1 路由前缀与分组

- **Decision**: 四渠道统一渠道前缀分组：`/common/*`、`/user/*`、`/shop/*`、`/admin/*`；组内资源路径。鉴权按组挂中间件（公开组/会员组/管理组）。
- **Rationale**: 140 端点无前缀必冲突（如会员与管理员的用户查询）；GoFrame 分组路由标准形态；现有 `/captcha`、`/login/sms` 迁移至 `/common/captcha`、`/user/login/sms`（003 契约的路径修订随本特性生效，见 D6）。
- **Alternatives**: 无前缀平铺（被拒：冲突且权限边界模糊）。

## D2 鉴权三级与承载

- **Decision**: 公开组无中间件；会员组挂 Bearer 会话校验（003 的 opaque token，含滑动续期）；管理组挂后台凭证 + 权限点校验（RBAC 权限树，`is_super` 直通）。管理端每个写接口登记权限点 `模块:资源:动作`（常量集于 `internal/consts`，与 V10 `admin_permission.code` 同构）。
- **Rationale**: 与既有 RBAC 表结构与 003 会话设计直接对接；权限点作为 RBAC 种子数据源。
- **Alternatives**: 注解式权限标注（被拒：Go 无注解，路由分组+常量最简）。

## D3 RefreshToken 纳管（用户实现反向对齐）

- **Decision**: 用户已在 `SmsLoginRes` 引入 `token + refreshToken` 双凭证——004 契约全线采用双凭证语义：token 短效（访问），refreshToken 长效（换新）；新增 `POST /user/token/refresh`。登出两者同失效。管理端同样双凭证。
- **Rationale**: 用户实现即事实契约；双凭证优于单滑动令牌的移动端体验（003 的"滑动续期"由 refresh 机制承担）。
- **Alternatives**: 回退单凭证（被拒：与已生成代码冲突）。

## D4 契约组织与文件粒度

- **Decision**: 一业务域一 v1 文件（如 `v1/order.go`、`v1/promotion.go`）；Req/Res 成对；渠道接口聚合文件（IUserV1 等）由 gf 工具维护标记为生成物；契约文档四渠道四文件 + conventions 总则。
- **Rationale**: 与用户既有生成形态一致；域文件避免单文件千行。

## D5 命名与字段约定

- **Decision**: 路径资源复数（`/user/addresses`）；JSON 字段 camelCase；金额字段单位元、类型 string（十进制安全，前端无精度丢失）；分页统一 `page/pageSize`，响应 `total/list`；ID 主键对外一律 string（int64 精度）；时间 RFC3339。
- **Rationale**: 既有 AGENTS 金额/精度约束的接口层落点；GoFrame 主流实践。

## D6 与 003 契约的衔接

- **Decision**: 003 的 8 端点路径按 D1 迁移（`/common/captcha*`、`/user/login/sms`、`/user/login/wx`、`/user/logout`、`/user/me`、新增 `/user/token/refresh`）；mock 取码端点保持条件注册；错误码分段延续并在本特性扩展（商品 3xxxx、交易 4xxxx、促销 5xxxx、分销 6xxxx、门店 7xxxx、后台 8xxxx 段）。
- **Rationale**: 单一演进方向，避免双契约；错误码段位在 003 D7 基础上按域扩展。
- **Alternatives**: 保留旧路径做兼容层（被拒：前端未开发，无兼容负担）。

## 关键事实核对

- api 目录四渠道已存在且 user/common 有 003 起步代码（gf 生成物形态）——契约与其同构 ✓
- errcode 包已存在（internal/errcode）——错误码扩展登记于此 ✓
- V10 admin_permission 表结构支持权限点种子 ✓
