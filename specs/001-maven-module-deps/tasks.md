---
description: "Task list for feature implementation"
---

# Tasks: Maven 模块依赖接线（分层依赖与版本仲裁）

**Input**: Design documents from `specs/001-maven-module-deps/`

**Prerequisites**: plan.md（必需）、spec.md（必需）、research.md、data-model.md、contracts/module-contracts.md、quickstart.md

**Tests**: 本特性验收以构建行为为证据（quickstart 四场景），未要求单元测试任务；FR-010 存量回归使用既有 `EcbootApplicationTests`。

**Organization**: 按用户故事分阶段：US1（构建通路）→ US2（版本仲裁）→ US3（方向强制）。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件，无未完成依赖）
- **[Story]**: 所属用户故事（US1/US2/US3）
- 描述中包含精确文件路径

## Path Conventions

- 本特性作用于仓库根的 Maven 模块：根 `pom.xml` 与 `dependencies/`、`infrastructure/`、`services/`、`apps/`、`start/` 下各 `pom.xml`
- 嵌套两层模块的父级相对路径为 `../../pom.xml`，一层为 `../pom.xml`

<!-- 由 /speckit-tasks 生成实际任务；下方各阶段任务为最终执行清单 -->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立前后对照证据基线

- [x] T001 记录改造前基线：在仓库根执行 `./mvnw validate`，将当前失败输出（父 POM 不可达形态）摘录到 `specs/001-maven-module-deps/baseline.md`，作为 US1 的对照证据

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 无——本特性无跨故事阻塞项；US1 本身即全部故事的基础增量（构建通路）。US2/US3 的任务隐式依赖 US1 完成后的可构建状态。

---

## Phase 3: User Story 1 - 一条命令构建全部后端模块 (Priority: P1) 🎯 MVP

**Goal**: 根 POM 成为聚合器 + 父级，10 个模块全部可解析父级并从仓库根一次构建成功

**Independent Test**: 全新环境（本地仓库无 org.juling 构件）在仓库根执行 `./mvnw clean package` → BUILD SUCCESS，reactor 顺序 infrastructure → services → apps → start，耗时 <1 分钟

### Implementation for User Story 1

- [x] T002 重写根 `pom.xml`：`<packaging>pom</packaging>`、parent 指向 `org.springframework.boot:spring-boot-starter-parent:4.1.1`、`<modules>` 声明 10 个模块（dependencies、infrastructure/ecboot-common、infrastructure/ecboot-infra-core、services/ecboot-service-user、services/ecboot-service-shop、apps/ecboot-api-common、apps/ecboot-api-user、apps/ecboot-api-shop、apps/ecboot-api-admin、start）、保留 `java.version=25`（依据 research.md D1/D2）
- [x] T003 [P] [US1] 修复 `dependencies/pom.xml`：parent 改为 `org.juling.ecboot:ecboot-parent:0.0.1-SNAPSHOT` 并使用相对路径 `../pom.xml`（删除阻断解析的 `<relativePath/>`）
- [x] T004 [P] [US1] 修复 `infrastructure/ecboot-common/pom.xml` 与 `infrastructure/ecboot-infra-core/pom.xml`：parent 指向 ecboot-parent，相对路径 `../../pom.xml`
- [x] T005 [P] [US1] 修复 `services/ecboot-service-user/pom.xml` 与 `services/ecboot-service-shop/pom.xml`：parent 指向 ecboot-parent，相对路径 `../../pom.xml`
- [x] T006 [P] [US1] 修复 `apps/ecboot-api-common/pom.xml`、`apps/ecboot-api-user/pom.xml`、`apps/ecboot-api-shop/pom.xml`、`apps/ecboot-api-admin/pom.xml`：parent 指向 ecboot-parent，相对路径 `../../pom.xml`
- [x] T007 修复 `start/pom.xml`：parent 由 spring-boot-starter-parent 改为 ecboot-parent，相对路径 `../../pom.xml`；依赖清单与插件声明本任务保持原样（框架版本来源由继承链获得，BOM 接线属 US2）
- [x] T008 [US1] 验证场景一：仓库根执行 `./mvnw clean package`，确认 BUILD SUCCESS、无 parent 解析错误、reactor 顺序符合分层、计时 <1 分钟；结果记入 `specs/001-maven-module-deps/quickstart.md` 验证记录区

**Checkpoint**: US1 完成——任何人在干净环境 clone 后一条命令可构建全部模块（MVP 可交付）

---

## Phase 4: User Story 2 - 版本号集中仲裁 (Priority: P1)

**Goal**: `dependencies` 成为依赖/框架版本唯一仲裁点；start 经继承链获取框架版本；插件版本集中于根 POM；业务 POM 零版本号

**Independent Test**: 检索全仓库 `<version>` 声明仅存在于 `dependencies/pom.xml`（属性+BOM 导入）与根 `pom.xml`（starter-parent 声明、enforcer 插件版本）；临时改 `spring-boot.version` 为无效值 → 构建失败；未登记依赖 → "version is missing" 失败（quickstart 场景二）

### Implementation for User Story 2

- [ ] T009 [US2] 在 `dependencies/pom.xml` 定义属性 `<spring-boot.version>4.1.1</spring-boot.version>`（框架版本全仓库唯一字面量），并在 `dependencyManagement` 导入 `org.springframework.boot:spring-boot-dependencies:${spring-boot.version}`（type=pom, scope=import）（依据 research.md D2）
- [ ] T010 [US2] 在根 `pom.xml` 的 `dependencyManagement` 导入 `org.juling:ecboot-dependencies:${project.version}`（type=pom, scope=import），使 BOM 仲裁经父级传递到全部模块（依据 research.md D2 就近优先规则）
- [ ] T011 [US2] 构建配置上移：将 `start/pom.xml` 中 maven-compiler-plugin 注解处理器路径（lombok + spring-boot-configuration-processor）、default-testCompile 处理器、hibernate-maven-plugin enhance 执行、native-maven-plugin 声明上移至根 `pom.xml` 的 build/pluginManagement；`start/pom.xml` 仅保留依赖清单与 spring-boot-maven-plugin 声明（依据 research.md D5/D6）
- [ ] T012 [US2] 验证场景二：从根 `./mvnw clean package` 复验构建成功；执行 quickstart 场景二三步（版本声明唯一性检索、spring-boot.version 无效值失败还原、ecboot-common 添加未登记依赖失败），结果记入 `specs/001-maven-module-deps/quickstart.md` 验证记录区

**Checkpoint**: US2 完成——版本单点改、全仓库生效；未登记/冲突版本在构建期失败

---

## Phase 5: User Story 3 - 依赖方向违规可检出 (Priority: P2)

**Goal**: enforcer bannedDependencies 按白名单矩阵（data-model.md）强制分层方向，绑定 validate 阶段；循环依赖构建失败

**Independent Test**: 临时注入 start→service 依赖 → 常规构建失败且含违规坐标；注入循环依赖 → reactor 报 cyclic reference；合规变更 → 构建通过（quickstart 场景三）

### Implementation for User Story 3

- [ ] T013 [US3] 在根 `pom.xml` 的 `pluginManagement` 声明 maven-enforcer-plugin（版本在此唯一声明，依据 research.md D6）：execution 绑定 `validate` 阶段，规则为 bannedDependencies（exclude 清单引用各模块属性 `${enforcer.banned.excludes}`，message 模板指明白名单契约文档路径）+ banDuplicatePomDependencyVersions + reactorModuleConvergence；并在根 build/plugins 激活，全部模块继承执行（依据 research.md D3，澄清 Q2 内置强制）
- [ ] T014 [P] [US3] 按契约（contracts/module-contracts.md）设置模块属性：`infrastructure/ecboot-common/pom.xml`（`org.juling:*`）、`infrastructure/ecboot-infra-core/pom.xml`（排除 ecboot-common）、`services/ecboot-service-user/pom.xml` 与 `services/ecboot-service-shop/pom.xml`（禁 `org.juling:ecboot`、`org.juling:ecboot-api-*`）
- [ ] T015 [P] [US3] 按契约设置 apps 层属性：`apps/ecboot-api-common/pom.xml`（禁 services/渠道/ecboot，允许 infrastructure——澄清 Q4）、`apps/ecboot-api-user/pom.xml`、`apps/ecboot-api-shop/pom.xml`、`apps/ecboot-api-admin/pom.xml`（各禁 start 与其他两个渠道）
- [ ] T016 [US3] 按契约设置 `start/pom.xml` 属性（禁 `org.juling:ecboot-service-*`、`org.juling:ecboot-infra-*`、`org.juling:ecboot-api-common`）与 `dependencies/pom.xml` 属性（禁全部 `org.juling:*` 业务构件）
- [ ] T017 [US3] 验证场景三与场景四：执行 quickstart 场景三全部四步注入（违规/循环/合规对照）确认失败与通过形态；执行 `cd start && ../mvnw test` 回归确认 FR-010，结果记入 `specs/001-maven-module-deps/quickstart.md` 验证记录区

**Checkpoint**: US3 完成——任何方向违规在任何一次常规构建中即时失败

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 证据归档与文档一致性

- [ ] T018 汇总归档：核对 `specs/001-maven-module-deps/quickstart.md` 验证记录区已含四个场景的命令、结果、耗时摘要（对照宪法 IV"以输出为证"）
- [ ] T019 文档同步：更新 `README.md`、`CLAUDE.md`、`AGENTS.md` 中"根 POM 非聚合器、需在 start/ 内构建"的构建说明与陷阱提示为新的"仓库根一条命令构建"事实，并移除已失效的 `mvnw -N install` 前置说明
- [ ] T020 终验：完整重跑 quickstart 四场景（干净本地仓库 org.juling 前提下），全部通过后本特性可提交

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (T001)**：无依赖，立即可做——产出 US1 的对照基线
- **US1 (T002–T008)**：依赖 T001 基线；T002 先行（根 POM 是全部 parent 解析的前提），T003–T006 可并行，T007 随后，T008 收口验证
- **US2 (T009–T012)**：依赖 US1 完成（可构建状态）；T009 → T010 → T011 → T012 顺序执行
- **US3 (T013–T017)**：依赖 US2 完成（enforcer 的 exclude 属性机制依赖根 pluginManagement 就位）；T013 先行，T014/T015 可并行，T016 随后，T017 收口验证
- **Polish (T018–T020)**：依赖 US3 完成

### User Story Dependencies

- US2、US3 均建立在 US1 的可构建状态之上（这是本特性内在的分层：先通路、再版本、再强制）
- US2 与 US3 之间：US3 依赖 US2 的根 pluginManagement 结构（T011），故顺序执行

### Parallel Opportunities

- T003–T006（US1 内，8 个模块 POM 两两独立）
- T014–T015（US3 内，6 个模块属性互不影响）
- 单人执行时按 ID 顺序即可；并行机会主要服务于多代理/多人场景

---

## Implementation Strategy

### MVP First (US1 Only)

1. T001 基线 → T002–T007 接线 → T008 验证
2. **STOP and VALIDATE**: 任何环境一条命令构建成功
3. 此时可交付：仓库具备可构建骨架（版本仲裁与强制尚未生效，但结构正确）

### Incremental Delivery

1. US1 → 构建通路 ✓（MVP）
2. +US2 → 版本单点仲裁 ✓
3. +US3 → 方向强制 ✓（完整交付 spec 全部验收标准）
4. Polish → 证据归档与文档一致性 ✓

---

## Notes

- 提交节奏：每个任务或逻辑任务组完成后提交一次，提交信息为中文 Conventional Commits（宪法 III），如 `feat: 根 POM 改造为聚合器并修复模块父级引用`
- 验证一律使用常规构建命令（`./mvnw validate` / `clean package` / `test`），以输出为证（宪法 IV）
- 违规注入类验证任务（T012/T017）完成后 MUST 还原临时修改
- 白名单矩阵的唯一权威来源：`specs/001-maven-module-deps/data-model.md` 与 `contracts/module-contracts.md`；实现时如有出入，修订契约并同步 spec，而非临场发挥
