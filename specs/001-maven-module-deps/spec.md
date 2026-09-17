# Feature Specification: Maven 模块依赖接线（分层依赖与版本仲裁）

**Feature Branch**: `001-maven-module-deps`

**Created**: 2026-09-17

**Status**: Draft

**Input**: User description: "设置pom.xml的依赖关系，参照最佳实践。dependencies为版本仲裁模块。start 模块仅依赖 apps模块；apps依赖 services模块；services模块依赖 infrastructure模块；"

## Clarifications

### Session 2026-09-17

- Q: start 模块的框架版本来源是否切换为经 dependencies 统一仲裁？ → A: 统一切换——start 放弃直接继承上游框架父级，经 ecboot-parent 继承链 / 导入 dependencies 物料清单获取框架版本（FR-005 保持严格，不设例外）
- Q: 分层/版本登记检查的强制时机？ → A: 内置强制——检查绑定到常规构建生命周期最早阶段，任何一次构建中违规或未登记版本即失败，无需独立命令或 CI
- Q: 构建插件版本是否纳入版本仲裁范围？ → A: 拆分治理——依赖/框架版本由 dependencies 唯一仲裁，构建插件版本集中在根父 POM 统一管理；业务模块 POM 零版本号（依赖或插件）
- Q: api-common 对 services 层的耦合强度？ → A: api-common MUST NOT 依赖 services 层；允许依赖 infrastructure 层——杜绝渠道模块经共享库传递引入领域实现

## User Scenarios & Testing *(mandatory)*

<!--
  用户故事按重要性排序（P1 最关键），每个故事必须可独立测试。
-->

### User Story 1 - 一条命令构建全部后端模块 (Priority: P1)

开发者在仓库根目录执行单条构建命令，构建系统按正确顺序
（infrastructure → services → apps → start）自动完成全部后端模块的构建，
全程无需任何手工前置步骤（例如手动安装父 POM 到本地仓库）。

**Why this priority**: 当前模块间没有可用的构建通路（父 POM 不可达、无聚合器），
任何依赖关系都无从验证；这是其余一切需求的前提。

**Independent Test**: 在干净环境 clone 仓库后，在根目录执行一次构建命令，
观察全部后端模块构建成功且日志中模块顺序符合分层方向。

**Acceptance Scenarios**:

1. **Given** 全新 clone 的仓库（本地仓库无本项目构件），**When** 开发者在仓库根执行一次构建，**Then** 全部后端模块构建成功，且无"找不到父 POM/构件"类错误
2. **Given** 仓库处于任意中间状态，**When** 执行根构建，**Then** 构建顺序遵循 infrastructure 先于 services、services 先于 apps、apps 先于 start
3. **Given** 单个模块构建失败，**When** 执行根构建，**Then** 构建在该模块失败并明确指出失败模块名

---

### User Story 2 - 版本号集中仲裁 (Priority: P1)

开发者需要升级或对齐某个第三方依赖版本时，只修改 `dependencies`
版本仲裁模块中的声明，重新构建后全仓库生效；业务模块的 POM 中不出现
任何具体版本号。

**Why this priority**: 版本仲裁是本次改造的核心目标之一；没有单一版本源，
依赖分层会立刻产生版本漂移与冲突。

**Independent Test**: 在代码库中检索所有版本声明位置，确认第三方与框架版本
仅存在于 `dependencies` 模块（及其上游框架清单导入处）；修改其中一处版本并
构建，验证全仓库模块解析到新版本。

**Acceptance Scenarios**:

1. **Given** 两个业务模块依赖同一第三方库，**When** 分别查看其 POM，**Then** 均未声明版本号，且构建产物解析到同一版本
2. **Given** `dependencies` 中声明的版本与上游框架清单默认版本不一致，**When** 构建业务模块，**Then** 以 `dependencies` 的声明为准
3. **Given** 业务模块开发者新增一个第三方依赖，**When** 未在 `dependencies` 登记版本，**Then** 常规构建失败，并提示应到 `dependencies` 声明

---

### User Story 3 - 依赖方向违规可检出 (Priority: P2)

当开发者（或后续 AI 代理）引入违反分层方向的依赖（例如 `services` 依赖
`apps`、`start` 直接依赖 `services`、出现循环依赖）时，常规构建必须失败，
并给出可定位的违规描述。

**Why this priority**: 分层规则只有"可被执行、可被检出"才有意义；但它建立在
US1 的可用构建与 US2 的版本源之上。

**Independent Test**: 在工作副本中临时加入一条违规依赖并执行常规构建，
确认构建失败且错误信息指明违规模块对；移除后恢复通过。

**Acceptance Scenarios**:

1. **Given** `start` 模块被加入对某个 `services` 模块的直接依赖，**When** 执行常规构建，**Then** 构建失败，错误信息指明 `start → services` 违规
2. **Given** 模块间被引入循环依赖（A → B → A），**When** 执行构建，**Then** 构建失败并报告循环路径
3. **Given** 合规的依赖变更（如 `services-user` 新增对 `infrastructure-common` 的依赖），**When** 执行常规构建，**Then** 构建通过

---

### Edge Cases

- 两个模块对同一第三方库传递引入了不同版本时，仲裁结果 MUST 确定、可解释（以 `dependencies` 声明为准），不允许"偶发地取决于声明顺序"
- `ecboot-api-common` 被渠道模块（api-user/api-shop/api-admin）依赖属于同层共享库，属合法依赖；若未来出现渠道模块互相依赖，MUST 判为违规
- `infrastructure` 层模块（ecboot-common、ecboot-infra-core）之间允许单向依赖（如 infra-core → common），MUST NOT 出现环
- 上游框架清单（如 Spring Boot 官方物料清单）与 `dependencies` 本地声明冲突时，MUST 有明确的优先级规则且全仓库一致

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `start` 模块 MUST 仅依赖 apps 层模块（ecboot-api-user、ecboot-api-shop、ecboot-api-admin），MUST NOT 直接依赖任何 `services/*` 或 `infrastructure/*` 模块
- **FR-002**: apps 渠道模块（api-user、api-shop、api-admin）MUST 通过 services 层模块获得业务能力；MUST NOT 依赖 `start`，MUST NOT 依赖其他渠道模块
- **FR-003**: services 模块（ecboot-service-user、ecboot-service-shop）MUST 仅依赖 infrastructure 层模块；MUST NOT 依赖 apps 层或 `start`
- **FR-004**: infrastructure 模块（ecboot-common、ecboot-infra-core）MUST NOT 依赖 services、apps 或 `start`；层内允许单向依赖
- **FR-005**: `dependencies` 模块 MUST 是全仓库第三方依赖与框架版本的唯一仲裁点；构建插件版本 MUST 集中在根父 POM 统一管理；业务模块 POM 中 MUST NOT 出现任何硬编码版本号（依赖或插件）；`start` MUST NOT 直接继承上游框架父级，框架版本 MUST 经由 `dependencies` 仲裁获取
- **FR-006**: 全部后端模块 MUST 可从仓库根以单条命令完成构建，且无需手动预装任何本仓库构件
- **FR-007**: 循环依赖 MUST 导致构建失败
- **FR-008**: 任何违反 FR-001～FR-004 分层方向的依赖声明 MUST 在常规构建中被自动检出（构建失败，检查绑定构建生命周期最早阶段），错误信息 MUST 包含违规的模块对
- **FR-009**: apps 层内共享模块 `ecboot-api-common` 可被同层渠道模块依赖；其自身 MUST NOT 依赖 services 层模块，MUST NOT 依赖渠道模块；允许依赖 infrastructure 层模块
- **FR-010**: 依赖仲裁规则变更后（版本升降、增删依赖），既有功能 MUST 无回归（现有测试全部通过，应用可正常启动）

### Key Entities *(include if feature involves data)*

- **ecboot-parent（仓库根 POM）**：全部模块的父级与聚合入口；定义 Java 版本等全局属性
- **ecboot-dependencies（版本仲裁模块）**：唯一允许声明第三方/框架版本的位置，以物料清单形式被其他模块消费
- **infrastructure 层**：ecboot-common、ecboot-infra-core——零业务语义的共享库，是依赖图的汇点
- **services 层**：ecboot-service-user、ecboot-service-shop——领域业务逻辑
- **apps 层**：ecboot-api-user、ecboot-api-shop、ecboot-api-admin（渠道 API）+ ecboot-api-common（同层共享）
- **start**：唯一可运行装配模块，组合 apps 暴露的能力
- **依赖关系（实体间关系）**：`start → apps → services → infrastructure` 严格单向；`dependencies` 被所有模块消费但不依赖任何业务模块

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 全新 clone 的仓库按 README 操作，首次构建即 100% 成功，无需任何口头指导或手工修复
- **SC-002**: 依赖与插件版本声明位置唯一——依赖/框架版本仅在 `dependencies` 模块，构建插件版本仅在根父 POM；版本调整单点修改即可全仓库生效
- **SC-003**: 构造任一方向违规依赖时，常规构建在构建阶段（而非运行期）失败，且错误信息可直接定位违规模块对
- **SC-004**: 改造完成后存量验证全部通过：现有测试通过、应用可正常启动，无功能回归
- **SC-005**: 从根目录完成一次全量构建的耗时可接受（全部模块当前为空实现，应在 1 分钟内完成）

## Assumptions

- "参照最佳实践"指采用 Maven 多模块的标准模式：根 POM 作为聚合器与父级
  （`packaging=pom` + 模块清单），版本仲裁通过 `dependencies` 模块的物料清单
  机制集中管理。当前"根 POM 非聚合器、子模块无法解析父级"的状态视为待修复问题。
- `ecboot-api-common` 定位为 apps 层内共享库（放渠道模块共用的 API 层代码），
  而非第四个渠道。
- 框架版本管理从 `start` 当前直接继承的上游框架父级统一切换至
  `dependencies` 仲裁（已澄清，见 Clarifications）；具体接线方式（父级链
  或清单导入）属实现设计，由 plan 阶段决定。
- 本特性不涉及前端项目，不新增/删除/重命名任何模块，不改变各模块既定职责。
- 分层/版本登记检查已澄清为内置强制：绑定常规构建生命周期最早阶段，
  任何一次构建即校验（见 Clarifications）；具体实现机制（构建插件、
  自定义脚本等）由 plan 阶段决定。
