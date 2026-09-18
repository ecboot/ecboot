# ecboot-start（装配模块）

唯一的**可运行模块**——把四渠道、领域服务与全部技术栈装配为单一 Spring Boot 应用（`EcbootApplication`）。

## 功能内容

| 职责 | 说明 |
|---|---|
| 应用入口 | `EcbootApplication`；`spring.application.name: ecboot` |
| 配置收敛 | `application.yaml`（数据源/Flyway 等）——运行时配置唯一落点（宪法 V） |
| 数据库迁移 | `src/main/resources/db/migration/`（Flyway 默认位置，V1~V22 / 54 表，随启动自动应用） |
| 渠道装配 | 依赖四渠道（api-user/shop/admin/common）——全部 REST API 的宿主 |
| 技术栈装配 | WebMVC、Security、JPA+Flyway、Redis、Elasticsearch、Quartz、Mail、WebSocket、RestClient、Validation、Actuator；运行时 MySQL 驱动、devtools、docker-compose 支持 |
| 构建产物 | spring-boot-maven-plugin repackage 可执行 jar；GraalVM native 支持（父 POM 声明） |

## 职责边界

- **做**：装配与配置；四渠道全装配（api-common 不装配则公共 API 不存在）
- **不做**：**不含业务代码**；禁直连 service/infra/common/webmvc（enforcer——必须经渠道传递）；不定义渠道内部行为

## 依赖关系

- 白名单：四渠道（api-common / api-user / api-shop / api-admin）
- 运行：`docker compose up -d`（MySQL 8.4 宿主 13306 / Redis / ES）→ `./mvnw spring-boot:run -pl ecboot-start -am`

相关文档：`../README.md`、`docs/schema-design.md`
