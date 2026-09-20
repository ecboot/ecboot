# API Contracts: 认证引导纵切片

统一响应契约（全部端点）：`{ "code": 0, "message": "ok", "data": ... }`；错误时 code≠0、message 用户可读、data 缺省。追踪号响应头 `X-Trace-Id`。

## 端点一览

### 1. GET `/api/captcha/image`（公开）

获取图形验证码。→ `data: { ticket, imageBase64, expiresIn }`（imageBase64 为 PNG data-URI）。

### 2. POST `/api/captcha/verify`（公开）

`{ ticket, answer }` → 校验图形码（一次性）。成功 `data:true`；失败 `code=20001`。

### 3. POST `/api/captcha/sms`（公开）

`{ phone, ticket, answer }`（先过图形码再发短信）→ `data:{ expiresIn }`。频控/锁定：`code=10004`（重发间隔）/`20002`（锁定）。

### 4. GET `/api/captcha/sms/mock-latest?phone=`（公开，**仅 mock 模式装配**）

返回 Redis 中的 mock 短信码。生产（mock 关闭）该路由不存在（404）。

### 5. POST `/api/auth/sms-login`（公开，api-user 渠道）

`{ phone, smsCode, channel? }` → `data: { token, userId, isNew, dormantVerified }`。错误：20001 验证码错/过期、20002 锁定、20003 禁用。

### 6. POST `/api/auth/wx-login`（公开）

`{ wxCode, phone?, smsCode?, channel? }`——mock 模式 phone 为显式开发字段。→ 同 5；归并冲突 `code=20004`（message 引导人工渠道）。

### 7. POST `/api/auth/logout`（受保护）

Header `Authorization: Bearer {token}`。→ 登出（token 即刻失效）。

### 8. GET `/api/auth/me`（受保护）

→ `data: { userId, nickname, avatar, registerChannel, lastLoginAt }`。未登录/失效 `code=10003`。

## 错误码契约（本切片）

| code | 语义 |
|---|---|
| 0 | 成功 |
| 10001 | 参数校验失败（message 含字段与原因） |
| 10002 | 系统错误（message 含追踪号） |
| 10003 | 未登录/凭证失效 |
| 10004 | 操作过于频繁（通用频控） |
| 20001 | 验证码错误或已过期 |
| 20002 | 尝试次数超限，已锁定 |
| 20003 | 账号已禁用 |
| 20004 | 微信身份与既有账号绑定冲突（人工渠道） |

## 安全契约

- 公共白名单集中于 `ecboot.security.public-paths`（start 配置）：验证码三端点 + 登录两端点 + actuator 健康。
- 受保护端点经 `AuthTokenFilter` 校验 Bearer token；`@CurrentUser` 注入 `userId`。
- 手机号明文只存在于请求/内存/日志脱敏输出；落库一律密文+哈希；响应体手机号脱敏 `138****1234`。
