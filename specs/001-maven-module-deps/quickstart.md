# Quickstart: 端到端验证手册

**Feature**: specs/001-maven-module-deps | **Date**: 2026-09-17

全部场景均以**常规构建命令的成败与输出为证据**（宪法 IV）。工作目录一律为仓库根
`D:\code\git\ecboot2`（命令用 `./mvnw`；PowerShell 下同样适用）。

## 前置条件

- JDK 25（GraalVM JDK 25 已装）
- Docker（仅场景四需要：用于 Spring Boot docker-compose 支持拉起 MySQL/Redis/ES）
- 全新验证前建议清理本地仓库中的旧构件：`rm -rf ~/.m2/repository/org/juling`

## 场景一：全新构建（US1 / FR-006 / SC-001 / SC-005）

```bash
./mvnw clean package
```

**期望**：
- `BUILD SUCCESS`，无任何 "Non-resolvable parent POM" / "Could not find artifact" 错误
- reactor 摘要显示 11 个项目全部成功，顺序满足 infrastructure → services → apps → start
- 计时（`time ./mvnw clean package`）在 1 分钟内（SC-005，空实现模块）

## 场景二：版本单点仲裁（US2 / FR-005 / SC-002）

1. 版本声明位置唯一：
   ```bash
   grep -rn "<version>" --include="pom.xml" . | grep -v "0.0.1-SNAPSHOT"
   ```
   **期望**：命中仅限 `dependencies/pom.xml`（spring-boot.version 属性与 BOM 导入）、
   根 `pom.xml`（enforcer 插件版本）、根继承链声明的 starter-parent 一处。
2. 单点修改生效（验证后还原）：
   - 临时把 `dependencies/pom.xml` 的 `spring-boot.version` 改为不存在的
     `99.99.99` → `./mvnw validate` **构建失败**，错误指向版本解析失败；
   - 还原为 `4.1.1` → `./mvnw validate` 恢复 `BUILD SUCCESS`。
3. 未登记版本即失败：向 `infrastructure/ecboot-common/pom.xml` 临时添加
   `org.apache.commons:commons-lang3`（不带版本）→ `./mvnw validate` 失败，
   报 `version is missing`；移除后恢复。

## 场景三：方向违规可检出（US3 / FR-007 / FR-008 / SC-003）

以下均使用临时修改 + `./mvnw validate`（最快阶段即触发），验证后还原：

1. `start/pom.xml` 临时添加 `org.juling:ecboot-service-user` 依赖 →
   **构建失败**，enforcer 输出含 `org.juling:ecboot-service-*` 违规说明；
2. `services/ecboot-service-user/pom.xml` 临时添加 `org.juling:ecboot-api-user` →
   同样失败；
3. `services/ecboot-service-user/pom.xml` 临时添加 `org.juling:ecboot-api-user` 且
   `apps/ecboot-api-user/pom.xml` 临时添加 `org.juling:ecboot-service-user` →
   构建失败，报告循环引用；
4. 合规对照：`services/ecboot-service-user/pom.xml` 临时添加
   `org.juling:ecboot-common` → 构建通过。

## 场景四：存量回归（FR-010 / SC-004）

```bash
cd start && ../mvnw test
```

**期望**：`EcbootApplicationTests` 通过（需要 Docker 拉起 compose 服务；本机无
Docker 时允许以 `-DskipTests` 验证纯打包路径并如实记录）。应用启动验证：
`../mvnw spring-boot:run` 正常启动无版本冲突告警（依赖收敛警告应为 0）。

## 验证记录

### 场景二（US2）：版本单点仲裁 — 2026-09-17 ✅

- 步骤 1 唯一性检索：`<version>` 声明仅存于 `dependencies/pom.xml`（spring-boot.version 属性 + BOM 导入）与根 `pom.xml`（starter-parent 继承链声明、BOM 导入 ${project.version}、hibernate 插件 ${hibernate.version}）✅ 业务模块零版本
- 步骤 2 单点生效：`spring-boot.version` 改为 99.99.99 → 构建失败 `Non-resolvable import POM: spring-boot-dependencies:pom:99.99.99`；还原 4.1.1 → BUILD SUCCESS ✅ 证明仲裁点在 `dependencies`
- 步骤 3 未登记依赖：**方法修正**——`validate` 阶段不解析依赖图，须用 `dependency:resolve`（或 compile）。探针一 commons-lang3 解析 3.20.0 成功：被 spring-boot-dependencies 托管，属"框架对齐=已登记"，仲裁链按设计生效；探针二 guava（框架未托管）→ `'dependencies.dependency.version' for com.google.guava:guava:jar is missing` ✅
- 还原后终验：`./mvnw clean package -DskipTests` → BUILD SUCCESS（8.8 s）
- 实现备注：BOM 模块独立化（不继承 ecboot-parent），否则根导入 BOM + BOM 继承根构成导入自环（Maven 报 "scope=import form a cycle"）；Maven 3.9.16 支持 reactor 内 import 解析

- `./mvnw clean package -DskipTests`（仓库根）：**BUILD SUCCESS**，11 个项目全部 SUCCESS
- reactor 顺序：ecboot-parent → dependencies → common → infra-core → service-user → service-shop → api-common → api-user → api-shop → api-admin → ecboot(start) ✅ 符合 infrastructure → services → apps → start
- 首次构建（含插件下载）1:03；热构建 **10.5 s**（SC-005 <1 分钟 ✅）
- 注：Docker 未运行，`-DskipTests` 为 quickstart 场景四允许的诚实回退；测试回归在 T017/T020 补验
- 对照基线：改造前根构建仅 1 项目、子模块 Non-resolvable parent POM（见 baseline.md）——RED→GREEN 闭环成立
