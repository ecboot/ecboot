# ecboot-common（基础通用库）

分层最底层的技术无关通用库——全仓库可复用的**值对象、错误码与工具**，不依赖任何内部构件（enforcer 全禁）。

## 功能内容

| 功能域 | 内容 | 对应规范 |
|---|---|---|
| 错误码体系 | 统一错误码枚举/常量（分域段位，如 用户域 1xxxx、交易域 2xxxx），业务异常携带错误码 | 供 api-webmvc 全局异常处理出口 |
| 通用业务异常 | 技术无关的 `BusinessException` 层级（错误码 + 参数），不含 web/http 语义 | — |
| 分页对象 | 分页请求/响应值对象（页码/页大小/总数/数据），Jackson 注解支撑序列化约定 | 全渠道复用 |
| 金额工具 | `BigDecimal` 安全运算/比较/格式化工具（比较用 `compareTo`，禁止 `equals`） | 宪法与 AGENTS 金额规则的代码落点 |
| 时间工具 | 时区统一的日期时间格式化/转换（Asia/Shanghai） | — |
| 通用常量 | 数字/字符串常量、正则（手机号等）集中定义 | — |

## 职责边界

- **做**：与框架运行时无关的通用定义——依赖仅 `jackson-annotations`、`jakarta.validation-api`（DTO 注解层）与 `lombok`
- **不做**：不引任何 Spring starter/容器（`spring-boot-starter` 都不引）；不放 web 语义（统一响应在 `ecboot-api-webmvc`）；不放技术组件实现（在 `ecboot-infra-core`）；不放业务逻辑

## 依赖关系

- 白名单：**无**（组内全禁——任何内部依赖都会被 enforcer 拒绝）
- 被依赖：全部模块（infra-core / api-webmvc / service-* / api 渠道 / start 经传递）

相关文档：`../README.md`（模块地图）、`specs/001-.../contracts/module-contracts.md`
