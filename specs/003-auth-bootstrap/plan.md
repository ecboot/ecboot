# Implementation Plan: 认证引导纵切片（auth-bootstrap）

**Branch**: `003-auth-bootstrap` | **Date**: 2026-09-20 | **Spec**: [spec.md](./spec.md)

## Summary

打通首个业务纵切片：统一响应/全局异常基座（common+api-webmvc）→ 图形/短信验证码组件（infra-core，短信 mock）→ 验证码端点（api-common）→ 手机号注册即登录/微信归并/休眠核身（service-user + api-user）→ 会话凭证（opaque token + Redis）。**零新表**（复用 user/user_login_log/system_config），验证六个模块的 enforcer 边界与 DDD 包结构首次实战。

## Technical Context

**Language/Version**: Java 25 + Spring Boot 4.1.1（继承链，零新框架）

**Primary Dependencies**: 既有栈（WebMVC/Security/Validation/JPA/Redis/Flyway）；**新增 1 个模块依赖**：infra-core 加 `spring-boot-starter-data-redis`（版本走 BOM，契约登记）——验证码存储/会话/频控均落 Redis

**Storage**: MySQL 8.4（零新迁移，V1~V31 已就绪）+ Redis（验证码凭证/会话/频控计数，键模型见 data-model.md）

**Testing**: `ecboot-start` 下 `@SpringBootTest` 端到端（TestRestTemplate 全链路：验证码→登录→受保护端点→登出）+ 组件单元测试（频控/一次性/加密往返）；宪法 IV 命令 `./mvnw test -pl ecboot-start -am`

**Target Platform**: apps/api 各模块（落位见 Project Structure）；运行依赖 compose 全栈

**Project Type**: 模块化单体纵切片（REST API 后端）

**Performance Goals**: 登录链路三接口 p95 < 500ms（本地容器口径）

**Constraints**: enforcer 矩阵（api-common 禁服务、渠道互禁、start 只装配渠道）；手机号库内零明文；全部阈值走 `system_config`（30 分钟生效）；B2C 严守（无 seller 交互面）

**Scale/Scope**: 新增约 25 个类、6 端点、0 迁移；6 模块首次落代码

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束 | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位既有模块（落位表见 research D6），无新模块；依赖方向经 enforcer（已验证矩阵） |
| II 统一技术栈 | ✅ | 零新第三方（验证码 AWT 自绘、token SecureRandom）；infra-core 加 starter-data-redis 属既有栈内补依赖，契约登记 |
| III 中文优先 | ✅ | 文档/注释中文；类名英文 |
| IV 可验证交付 | ✅ | 端到端测试 + quickstart 可复制命令（curl/PowerShell 序列）；mock 短信码经 Redis 查询端点获取（仅 mock 开关暴露） |
| V 简单优先 | ✅ | 零新表、零新库；opaque token 而非 JWT（YAGNI）；会话/验证码收敛 Redis |
| 工程约束 1.2.0 | ✅ | 迁移位置/构建目录合规；**多商户预留（B2C 严守）**：本切片为用户域，无 seller 交互面，无需维度注入（显式声明） |

**Phase 1 复查**：✅ 设计产物（research/data-model/contracts/quickstart）未越界——无新模块、无新第三方、无 seller 功能泄漏。

## Project Structure

### Documentation (this feature)

```text
specs/003-auth-bootstrap/
├── plan.md / research.md / data-model.md / quickstart.md
├── contracts/api-contracts.md
└── tasks.md               # /speckit-tasks 生成
```

### Source Code (repository root)

```text
apps/api/
├── ecboot-common/src/main/java/org/juling/ecboot/common/
│   ├── api/ApiResponse.java, ErrorCode.java, PageResult.java
│   └── exception/BusinessException.java
├── ecboot-infra-core/src/main/java/org/juling/ecboot/infra/
│   ├── captcha/CaptchaService.java, CaptchaTicket.java, ImageCaptchaGenerator.java
│   ├── sms/SmsSender.java, MockSmsSender.java
│   ├── security/PhoneCipher.java, SessionManager.java, SessionKeys.java
│   └── config/InfraCoreAutoConfiguration.java, RedisConfig.java
├── ecboot-api-webmvc/src/main/java/org/juling/ecboot/web/
│   ├── GlobalExceptionHandler.java, TraceIdFilter.java, WebJacksonConfig.java
│   ├── security/AuthTokenFilter.java, CurrentUser.java, CurrentUserResolver.java,
│   │   SecurityBaselineConfig.java（读 ecboot.security.public-paths 白名单）
│   └── config/WebAutoConfiguration.java
├── ecboot-api-common/src/main/java/org/juling/ecboot/apicommon/
│   └── captcha/CaptchaController.java（image/verify/sms + mock 码查询端点）
├── ecboot-service-user/src/main/java/org/juling/ecboot/user/
│   ├── interfaces/dto/SmsLoginCommand.java, WxLoginCommand.java, LoginResult.java ...
│   ├── application/AuthAppService.java（注册即登录/归并/休眠核身）
│   ├── domain/model/User.java, UserLoginLog.java（JPA 实体，映射既有表）
│   ├── domain/repo/UserRepository.java, UserLoginLogRepository.java
│   └── infrastructure/persistence/（JPA 实现）
├── ecboot-api-user/src/main/java/org/juling/ecboot/apiuser/
│   └── auth/AuthController.java（sms-login/wx-login/logout/me）
└── ecboot-start/src/main/resources/application.yaml（phone-key/公共白名单/mock 开关）
```

**Structure Decision**: DDD 四层包结构在 service-user 首次落地（interfaces/application/domain/infrastructure）；渠道只做 DTO 组装；跨层数据经 interfaces 包的应用接口。

## Complexity Tracking

> 无宪法违例。最接近边界的一项：mock 短信码查询端点（仅开发态条件装配）——属测试便利设施，非生产功能，已在 contracts 标注暴露条件。
