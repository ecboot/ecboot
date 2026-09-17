# Implementation Plan: Maven 模块依赖接线（分层依赖与版本仲裁）

**Branch**: `001-maven-module-deps` | **Date**: 2026-09-17 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-maven-module-deps/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command; its definition describes the execution workflow.

## Summary

将仓库根 POM 改造为聚合器 + 父级，建立 `dependencies` 唯一版本仲裁点（导入
Spring Boot 官方清单、集中第三方版本），以 `ecboot-parent → spring-boot-starter-parent`
继承链为 `start` 提供框架版本管理，并通过 maven-enforcer 的 bannedDependencies
规则（绑定常规构建最早阶段）机械强制 `start → apps → services → infrastructure`
的单向依赖与禁止循环，使全部后端模块可从仓库根一条命令构建成功。

## Technical Context

**Language/Version**: Java 25（GraalVM JDK 25，已安装）；Maven 3.9.16（仓库自带 wrapper `mvnw`）

**Primary Dependencies**: Maven 原生多模块机制（聚合器、父 POM 继承、dependencyManagement 物料清单导入、pluginManagement）；maven-enforcer-plugin（依赖方向强制）；Spring Boot 4.1.1（start 现有 starters，仅作为被仲裁对象，不新增 starter）

**Storage**: N/A（本特性为构建体系改造，无数据存储；compose.yaml 基础设施不在范围内）

**Testing**: 以 Maven 构建行为为验收手段——常规构建成败即测试（`../mvnw clean verify` 从仓库根）；存量回归用 `start/` 下 `../mvnw test`（EcbootApplicationTests 上下文加载测试）

**Target Platform**: JVM（本地开发机，Windows 11 + Git Bash/PowerShell；不引入 CI）

**Project Type**: Maven 多模块单体（11 个 POM：1 根聚合器 + 10 模块）

**Performance Goals**: 全量构建（当前全部模块为空实现）在 1 分钟内完成（SC-005）

**Constraints**: 不引入 CI 与新增顶层模块；不使用前端项目；`dependencies` 为依赖/框架版本唯一仲裁点、根父 POM 为插件版本唯一集中点（澄清 Q1–Q4 已锁定）

**Scale/Scope**: 修改 11 个 POM + 新增 4 份设计文档；不新增/删除/重命名模块，不触碰前端与业务代码

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则 | 门禁问题 | 状态 | 依据 |
| --- | --- | --- | --- |
| I. 模块化单体 | 依赖方向是否严格单向 start → apps → services → infrastructure，无循环？ | ✅ PASS | 本特性即该原则的落地机制：bannedDependencies 按层白名单 + reactor 天然循环失败（FR-001～004、007、008） |
| II. 统一技术栈 | 是否混入同类替代框架/运行时？ | ✅ PASS | 仅使用 Maven 原生机制与 maven-enforcer（构建工具，非应用框架）；不新增 starter/语言/前端依赖 |
| III. 中文优先 | 文档与提交信息是否中文？标识符英文？ | ✅ PASS | 本计划及 research/data-model/contracts/quickstart 均中文；模块坐标保持英文 |
| IV. 可验证交付 | 每步是否以构建输出为证？ | ✅ PASS | quickstart 的全部验证场景均以常规构建命令成败与输出为证据；无"应已通过"式断言 |
| V. 简单优先 | 是否新增顶层模块/抽象层？机制是否最小？ | ✅ PASS | 不新增模块；方向检查用 enforcer 单插件 + 根级配置参数化；不引入第三方 enforcer 扩展或自研检查脚本；插件版本集中仅新增 enforcer 一项 |

**结论**: 无违例，`Complexity Tracking` 表留空。

## Project Structure

### Documentation (this feature)

```text
specs/001-maven-module-deps/
├── plan.md              # 本文件（/speckit-plan 输出）
├── research.md          # Phase 0 输出：机制选型与决策记录
├── data-model.md        # Phase 1 输出：模块实体、依赖白名单矩阵、版本仲裁数据流
├── contracts/           # Phase 1 输出：各 POM 的坐标/继承/规则契约
│   └── module-contracts.md
├── quickstart.md        # Phase 1 输出：端到端验证手册
└── tasks.md             # Phase 2 输出（/speckit-tasks 生成，非本命令产物）
```

### Source Code (repository root)

```text
pom.xml                              # 根：ecboot-parent —— 改为 packaging=pom 聚合器 + 父级
dependencies/
└── pom.xml                          # ecboot-dependencies —— 版本仲裁 BOM（依赖/框架版本唯一声明点）
infrastructure/
├── ecboot-common/pom.xml            # 基础库，依赖白名单：空（不依赖业务模块）
└── ecboot-infra-core/pom.xml        # 基础库，允许：ecboot-common
services/
├── ecboot-service-user/pom.xml      # 允许：infrastructure/*
└── ecboot-service-shop/pom.xml      # 允许：infrastructure/*
apps/
├── ecboot-api-common/pom.xml        # 允许：infrastructure/*；禁止：services/*、渠道模块
├── ecboot-api-user/pom.xml          # 允许：services/*、api-common、infrastructure/*；禁止：start、其他渠道
├── ecboot-api-shop/pom.xml          # 同上
└── ecboot-api-admin/pom.xml         # 同上
start/
└── pom.xml                          # parent 改为 ecboot-parent；允许：apps 渠道模块；禁止：services/*、infrastructure/*
```

**Structure Decision**: 沿用仓库既有 11 模块骨架，仅重写各 POM 的 parent/依赖/
规则配置，不改变目录布局；分层依赖白名单的完整矩阵与每模块 enforcer 参数见
[data-model.md](./data-model.md) 与 [contracts/module-contracts.md](./contracts/module-contracts.md)。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| （无） | — | — |
