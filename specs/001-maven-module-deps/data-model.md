# Data Model: Maven 模块依赖接线

**Feature**: specs/001-maven-module-deps | **Date**: 2026-09-17

本特性为构建体系改造，"领域数据"即：模块实体、依赖边（白名单矩阵）、版本仲裁流。
无生命周期状态迁移；构建行为验证规则见文末映射表。

## 实体：POM 模块（14 个）

| 构件（GA: org.juling.ecboot:*） | packaging | parent | 角色 |
| --- | --- | --- | --- |
| ecboot-parent | pom | spring-boot-starter-parent (外部) | 聚合器 + 根父级 + 插件版本集中点 |
| ecboot-dependencies | pom | 无（独立 BOM，避免导入自环） | 依赖/框架版本唯一仲裁 BOM |
| ecboot-infrastructure | pom | ecboot-parent | infrastructure 层聚合器（用户决策 2026-09-17 接入） |
| ecboot-services | pom | ecboot-parent | services 层聚合器（同上） |
| ecboot-apps | pom | ecboot-parent | apps 层聚合器（同上） |
| ecboot-common | jar | ecboot-parent | infrastructure 基础库（汇点） |
| ecboot-infra-core | jar | ecboot-parent | infrastructure 基础库 |
| ecboot-service-user | jar | ecboot-parent | user 领域服务 |
| ecboot-service-shop | jar | ecboot-parent | shop 领域服务 |
| ecboot-api-common | jar | ecboot-parent | 公共渠道 API（如短信，图片验证码等） |
| ecboot-api-user | jar | ecboot-parent | user 渠道 API |
| ecboot-api-shop | jar | ecboot-parent | shop 渠道 API |
| ecboot-api-admin | jar | ecboot-parent | admin 渠道 API |
| ecboot | jar | ecboot-parent | 唯一可运行装配模块（start） |

根 `<modules>` 聚合 5 个条目：`dependencies`（BOM 须经根 reactor 供导入解析）、
三个层聚合器、`start`；层聚合器各自聚合层内叶子模块（聚合器 ≠ 父级，叶子
parent 仍为 ecboot-parent）。

## 关系：依赖边白名单矩阵（enforcer bannedDependencies 的依据）

"允许"指对 org.juling.ecboot 业务构件的白名单；第三方构件不受方向规则约束（仅受版本纪律约束）。未列入"允许"的业务构件一律进入该模块的 enforcer 禁止清单。

| 模块 | 允许依赖的业务构件 | 明确禁止 |
| --- | --- | --- |
| ecboot-dependencies | 无 | 任何业务构件（BOM 不依赖业务） |
| ecboot-common | 无 | 任何业务构件 |
| ecboot-infra-core | ecboot-common | 其余全部业务构件 |
| ecboot-service-user / shop | ecboot-common、ecboot-infra-core | ecboot、全部 ecboot-api-*、同层兄弟 service |
| ecboot-api-common | ecboot-common、ecboot-infra-core | ecboot-service-*、全部渠道 ecboot-api-*（澄清 Q4：禁依赖 services）、ecboot |
| ecboot-api-user / shop / admin | ecboot-api-common、ecboot-service-*、ecboot-common、ecboot-infra-core | ecboot、其他两个渠道模块 |
| ecboot (start) | ecboot-api-user、ecboot-api-shop、ecboot-api-admin | ecboot-service-*、ecboot-infra-*、ecboot-api-common（FR-001 白名单枚举之外） |

矩阵性质：严格无环；任意模块到 ecboot-common 存在唯一方向的有向路径；`start` 到任意 services/infrastructure 构件不存在直接边。

## 版本仲裁数据流

```text
[字面量 4.1.1] spring-boot.version ── 唯一声明于 ecboot-dependencies
        │
        ├─► ecboot-dependencies 导入 spring-boot-dependencies（框架依赖版本）
        │        ▲ import
ecboot-parent dependencyManagement ──► 传递给全部子模块（就近优先，覆盖祖父级 starter-parent 默认值）
        │
        ├─► 第三方版本（未来）声明于 ecboot-dependencies dependencyManagement
        │
模块 POM：零 <version>（依赖与插件）——未登记即"version missing"构建失败
```

版本字面量全仓库仅两处：`dependencies/pom.xml` 的 `spring-boot.version` 属性、
根 `pom.xml` 继承链声明的 starter-parent 版本（4.1.1）。enforcer 插件版本由
starter-parent 托管（评审修订，原显式 3.5.0 已移除）。模块自身版本
`0.0.1-SNAPSHOT` 由 parent 坐标定义，子模块经继承获得。

## 验证规则映射（FR → 机制）

| 需求 | 机制 | 失败表现 |
| --- | --- | --- |
| FR-001～004（方向） | enforcer bannedDependencies（根级配置 + 模块属性参数化，validate 阶段） | 构建失败，输出含违规构件坐标 |
| FR-005（版本唯一） | BOM 导入 + 插件版本集中于根 pluginManagement + banDuplicatePomDependencyVersions | 缺版本即"version missing"失败；重复声明失败 |
| FR-006（根构建） | 聚合器 + 相对路径父级（D1） | —（成功路径） |
| FR-007（禁循环） | reactor 原生拓扑检测 + 白名单结构性无环 | "cyclic reference" 构建失败 |
| FR-008（违规可检出） | 同 FR-001～004，validate 阶段内置 | 同上 |
| FR-009（api-common 耦合） | 矩阵中 api-common 行的禁止清单 | 同 enforcer |
| FR-010（无回归） | 存量 `start/` 下 `../mvnw test` | 测试通过即合规 |
