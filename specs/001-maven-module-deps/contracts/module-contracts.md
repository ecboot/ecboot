# Contracts: 模块 POM 契约

**Feature**: specs/001-maven-module-deps | **Date**: 2026-09-17

每个 POM 对外暴露的稳定契约：坐标、继承、packaging、业务依赖白名单（enforcer
属性）、对外提供物。实现（tasks 阶段）与评审以此为准；违反即视为破坏契约。

## 根：ecboot-parent（pom.xml）

- 坐标 `org.juling.ecboot:ecboot-parent:0.0.1-SNAPSHOT`，`packaging=pom`
- parent：`org.springframework.boot:spring-boot-starter-parent:4.1.1`
- 聚合 modules（5，用户决策 2026-09-17：层聚合器间接聚合叶子模块）：
  `dependencies`（BOM 必须经根 reactor 供导入解析）、`infrastructure`、
  `services`、`apps`、`start`
- properties：`java.version=25`（保留），新增 enforcer 参数默认值
- dependencyManagement：import `org.juling.ecboot:ecboot-dependencies:${project.version}`
- pluginManagement / build：compiler（lombok + configuration-processor 注解处理器、
  release 25）、hibernate 增强（自 start 上移）、maven-enforcer（bannedDependencies
  **searchTransitive=false 仅查直接依赖** + banDuplicatePomDependencyVersions，绑定
  validate）；enforcer 版本由 starter-parent 继承链托管（实现修订：ReactorModuleConvergence
  因独立 BOM 移除，enforcer 显式版本移除）
- **对全仓库承诺**：从根目录 `mvnw` 任意生命周期一次构建全部模块

## BOM：ecboot-dependencies（dependencies/pom.xml）

- 坐标 `org.juling.ecboot:ecboot-dependencies:0.0.1-SNAPSHOT`，`packaging=pom`
- parent：**无（独立 BOM）**——实现修订：若继承 ecboot-parent，"根导入 BOM + BOM 继承根"会构成 import 自环
- properties：`spring-boot.version=4.1.1`（框架版本全仓库唯一字面量）
- dependencyManagement：import `spring-boot-dependencies:${spring-boot.version}`；
  未来第三方版本仅在此追加
- 禁止依赖任何 org.juling 业务构件

## 层聚合器（infrastructure/、services/、apps/ 的 pom.xml）

- 坐标 `org.juling.ecboot:ecboot-infrastructure|ecboot-services|ecboot-apps:0.0.1-SNAPSHOT`，`packaging=pom`
- parent：ecboot-parent（`../pom.xml`）；`<modules>` 聚合层内叶子模块——
  **聚合器 ≠ 父级**，叶子 parent 仍为 ecboot-parent（`../../pom.xml`）
- MUST NOT 声明任何依赖；enforcer 默认空禁令经继承覆盖（packaging=pom 无影响）
- `dependencies` BOM 不属于任何层，独立挂于根 modules（其 import 由根 reactor 解析）

## infrastructure 层

| 模块 | 业务依赖白名单 | enforcer 禁止清单（属性值，逗号分隔 GA 通配） |
| --- | --- | --- |
| ecboot-common | 无 | `org.juling.ecboot:*` |
| ecboot-infra-core | ecboot-common | `org.juling.ecboot:*` 排除 `org.juling.ecboot:ecboot-common` |

## services 层

| 模块 | 业务依赖白名单 | enforcer 禁止清单 |
| --- | --- | --- |
| ecboot-service-user | ecboot-common、ecboot-infra-core | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-api-*`、`org.juling.ecboot:ecboot-service-shop`（同层互禁，评审修订） |
| ecboot-service-shop | 同上（兄弟禁令互换） | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-api-*`、`org.juling.ecboot:ecboot-service-user` |

## apps 层

| 模块 | 业务依赖白名单 | enforcer 禁止清单 |
| --- | --- | --- |
| ecboot-api-common | ecboot-common、ecboot-infra-core | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-service-*`、`org.juling.ecboot:ecboot-api-user`、`org.juling.ecboot:ecboot-api-shop`、`org.juling.ecboot:ecboot-api-admin` |
| ecboot-api-user | ecboot-api-common、ecboot-service-*、ecboot-common、ecboot-infra-core | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-api-shop`、`org.juling.ecboot:ecboot-api-admin` |
| ecboot-api-shop | 同上（渠道白名单互换） | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-api-user`、`org.juling.ecboot:ecboot-api-admin` |
| ecboot-api-admin | 同上（渠道白名单互换） | `org.juling.ecboot:ecboot`、`org.juling.ecboot:ecboot-api-user`、`org.juling.ecboot:ecboot-api-shop` |

## 装配：ecboot（start/pom.xml）

- 坐标 `org.juling.ecboot:ecboot:0.0.1-SNAPSHOT`，`packaging=jar`
- parent：**改为 ecboot-parent**（不再直接继承 spring-boot-starter-parent）
- 业务依赖：ecboot-api-user、ecboot-api-shop、ecboot-api-admin
- 第三方依赖：现有 12 个 Boot starter + flyway-mysql、mysql-connector-j、
  devtools、docker-compose、lombok、12 个 test starter —— **全部保留**，
  无版本号（经继承链管理）
- enforcer 禁止清单：`org.juling.ecboot:ecboot-service-*`、`org.juling.ecboot:ecboot-infra-*`、
  `org.juling.ecboot:ecboot-api-common`
- build：保留 spring-boot-maven-plugin 声明（repackage 由 starter-parent 托管）；
  移除已上移的 compiler/hibernate 配置

## 消费契约（模块使用方视角）

1. 任何模块新增第三方依赖：**仅在 `dependencies/pom.xml` 登记版本**，模块内零版本引用；
2. 任何模块新增业务依赖：仅允许引用白名单内构件，白名单变更需修订本契约与 spec；
3. 新增插件版本：仅根 `pluginManagement`；
4. 新增模块：挂入根 modules、parent 指向 ecboot-parent、按矩阵补 enforcer 属性，
   并同步更新 data-model 矩阵。

## 扁平结构落地修订（2026-09-18，POM 重构）

仓库重构后模块扁平挂于 `apps/api/`（层聚合器已删除），本节为权威矩阵的现行形态：

- **禁令机制**：逐 GA 精确列举，槽位 `enforcer.banned.1..8`（默认 `__none__`）；
  **`allowed` includes 机制已移除**。实证（enforcer 3.6.3）：`bannedDependencies`
  仅支持精确 `groupId:artifactId` 匹配——纯 groupId 与 `groupId:*` 通配**均无效**
  （2026-09-18 注入实验：exclude=`org.juling.ecboot` 时 api-user←api-shop 通过）。
- **禁令 = 内部构件全集（9）− 本模块白名单**；反向依赖另由 reactor 环检测兜底拦截。
- 白名单矩阵：common=∅；infra-core={common}；service-*= {common, infra-core}（兄弟互禁）；
  api-common={common, infra-core}（全部 service 禁入）；api-user/shop={api-common,
  service-*, common, infra-core}（兄弟渠道互禁）；api-admin 同左；start={api-user,
  api-shop, api-admin}（service/infra/common 直连禁入，必须经渠道传递）。
- **各层技术依赖登记**（零版本号，版本经 BOM 继承链）：
  common=jackson-annotations、jakarta.validation-api、lombok(optional)；
  infra-core=spring-boot-starter、configuration-processor(optional)；
  service-*=data-jpa、validation、lombok(optional)、starter-test(test)；
  api-common=webmvc、validation、lombok(optional)；
  api-user/shop=api-common+对应 service+starter-test(test)；api-admin=api-common+双 service+test；
  start=依赖清单不变（12 starter + flyway-mysql + mysql 驱动 + devtools/docker-compose + lombok + test starters）。
- 模块 POM 不再冗余声明 `java.version`（继承父级）。
