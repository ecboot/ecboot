# ecboot-api-common（公共 API 渠道）

与 user/shop/admin **同级的第四渠道**——对外提供无业务归属的**公共 REST API**，当前以短信与图形验证码为主。

## 功能内容（REST 端点规划）

| 端点（规划） | 功能 | 依赖组件 |
|---|---|---|
| `POST /api/captcha/image` | 生成图形验证码（返回图片+临时凭证），配套校验接口 | infra-core 验证码组件 |
| `POST /api/captcha/sms` | 短信验证码下发（手机号+图形验证码前置校验；**频控与防刷**：单号/单 IP 限频） | infra-core 短信/验证码组件 |
| `POST /api/captcha/sms/verify` | 短信验证码校验（注册/登录前置） | 同上 |
| `GET /api/ping` 等探针 | 公共健康/联通探针（业务无关） | — |

> 防刷与风控联动：下发异常行为命中风控规则时联动拦截（风控域在 service-shop，经事件/接口协作）。

## 职责边界

- **做**：公共能力 的 **REST 暴露层**——薄控制器，实现在 infra-core 组件
- **不做**：**禁依赖任何领域服务**（enforcer，宪法要求）；与 user/shop/admin 互不编译依赖（四渠道互禁——其他渠道需要发码能力时复用 infra-core 组件，不依赖本渠道）

## 依赖关系

- 白名单：`ecboot-api-webmvc`（Web 基座）、`ecboot-common`、`ecboot-infra-core`
- 被依赖：仅 `ecboot-start`（四渠道全装配，公共 API 随应用可用）

相关文档：`../README.md`、`ecboot-infra-core/README.md`（验证码/短信组件）
