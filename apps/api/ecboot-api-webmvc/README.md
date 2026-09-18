# ecboot-api-webmvc（Web 基座）

四渠道（api-common/user/shop/admin）共同依赖的 **Web 技术基座**——统一响应、异常、安全与通用 Web 配置的唯一定义处。

## 功能内容

| 功能域 | 内容 |
|---|---|
| 统一响应 | `Result<T>` 包装模型（code/message/data）、成功/失败构造、错误码出口约定（码表在 ecboot-common） |
| 全局异常处理 | 业务异常→错误码响应；参数校验异常→字段级 400 明细；未知异常→500 + 追踪号（不泄漏堆栈） |
| Web 通用配置 | CORS 白名单、全局序列化规则（日期格式/时区、Long→String 防前端精度丢失）、统一编码 |
| 安全基线 | Spring Security 过滤器链统一装配（路径白名单、认证入口、登录态解析），各渠道只声明路径权限 |
| 登录态 | 当前用户上下文（`@CurrentUser` 参数解析器：从会话/JWT 解析会员或后台账号身份） |
| 请求日志 | 入口日志与追踪号（TraceId 贯穿日志/响应头） |
| 通用注解/切面 | `@RequireLogin` 等渠道可复用的声明式守卫 |

## 职责边界

- **做**：与"Web 如何呈现"相关的全部横切配置；技术依赖 = `spring-boot-starter-webmvc` + `validation`（四渠道经本模块传递，**不再直接声明 web starter**）
- **不做**：**纯技术基座，禁依赖任何领域服务与渠道**（enforcer 强制）；无业务端点；业务异常码表在 ecboot-common

## 依赖关系

- 白名单：`ecboot-common`、`ecboot-infra-core`
- 被依赖：api-common / api-user / api-shop / api-admin（四渠道的共同基座）

相关文档：`../README.md`（模块地图）、`specs/001-.../contracts/module-contracts.md` §矩阵修订 v2
