# ECBOOT 后端（apps/api）

社交电商平台后端——Java 25 / Spring Boot 4.1.1 的**模块化单体**，由本目录（`ecboot-parent`，Maven 聚合器 + 父 POM）统一聚合，`ecboot-start` 装配为单一部署单元。

## 模块地图

```
ecboot-start ── 装配四渠道（REST API 的唯一宿主）
├─ ecboot-api-user    小程序前台 · 会员中心 API
├─ ecboot-api-shop    小程序前台 · 商城 API
├─ ecboot-api-admin   管理后台 API
├─ ecboot-api-common  公共 API（短信/图形验证码等）
│     └── 四渠道共同依赖 → ecboot-api-webmvc（Web 基座：统一响应/异常/安全基线）
│                          └─→ ecboot-common / ecboot-infra-core（基础层）
└─ ecboot-service-user / ecboot-service-shop（领域服务，内部按 DDD 分层）
ecboot-dependencies（独立 BOM：全仓库版本仲裁唯一字面量）
```

各模块有独立 README（定位/职责边界/功能内容/依赖白名单）。

## 分层与强制规则

- 依赖方向：`start → api 渠道 → service 领域 → common/infra-core`，**四渠道互不编译依赖**，service 兄弟互禁——全部由 maven-enforcer 在 `validate` 阶段强制（逐 GA 精确禁令矩阵），违规即构建失败。矩阵权威：`specs/001-maven-module-deps/contracts/module-contracts.md`。
- 版本仲裁：第三方/框架版本仅在 `ecboot-dependencies` 声明；插件版本仅在根 `pluginManagement`；业务模块 POM 零版本号。

## 构建与运行

```bash
cd apps/api
docker compose up -d                 # MySQL 8.4(宿主13306) / Redis / Elasticsearch
./mvnw clean package                 # 全 reactor 构建（含测试）
./mvnw test -pl ecboot-start -am     # 宪法 IV 验证命令
./mvnw spring-boot:run -pl ecboot-start -am   # 启动（Flyway 自动应用 db/migration 全部迁移）
```

## 数据库

- 迁移：`ecboot-start/src/main/resources/db/migration/`（Flyway 默认位置，V1~V22，54 张表）
- Schema 设计文档：`docs/schema-design.md`；验证指南：`specs/002-social-commerce-expansion/quickstart.md`

## 关键约束（宪法摘要）

- 中文 Conventional Commits；文档中文、标识符英文
- 金额一律 `DECIMAL(10,2)` + `BigDecimal.compareTo`，禁止 float/equals
- 合规红线：分销≤2 级（结构强制）、注销匿名化、敏感字段加密+哈希索引
- 完整宪法：`.specify/memory/constitution.md`
