# Research: Maven 模块依赖接线

**Feature**: specs/001-maven-module-deps | **Date**: 2026-09-17

本文件记录 Phase 0 的机制选型决策。所有决策以 spec 的 FR/澄清（Q1–Q4）为约束。

## D1: 根 POM 改造为聚合器 + 父级

**Decision**: `pom.xml`（ecboot-parent）改为 `<packaging>pom</packaging>` 并声明
`<modules>` 清单（dependencies、infrastructure/*、services/*、apps/*、start 共 10 个
模块）；全部子模块的 `<parent>` 使用相对路径引用（嵌套两层目录用 `../../pom.xml`），
删除现有阻断解析的 `<relativePath/>` 空标记。

**Rationale**: Maven 多模块标准形态——reactor 自动按依赖拓扑排序构建（满足 FR-006、
US1），子模块可直接从文件系统解析父级（消除"父 POM 不可达"断路）；`<modules>`
顺序无关，构建顺序由 reactor 计算。

**Alternatives considered**:
- 保持根 POM 非聚合器、各模块独立构建 → 构建顺序靠手工 install，直接违反 FR-006/US1。
- 引入独立的 aggregator POM（parent 与 aggregator 分离）→ 多一个顶层文件，当前规模
  无收益，违反宪法 V（简单优先）。

## D2: 父级链与 BOM 优先级（澄清 Q1 的接线机制）

**Decision**:

```text
spring-boot-starter-parent (4.1.1, 外部)
        ↑ 继承
ecboot-parent (根: 聚合器 + pluginManagement + dependencyManagement 导入 BOM)
        ↑ 继承                    ↑ import（依赖管理导入）
全部 10 个模块              ecboot-dependencies (BOM)
        ↑ 继承                     ↑ import
start ──┘                  spring-boot-dependencies (${spring-boot.version} 属性)
```

- `ecboot-dependencies`：定义 `<spring-boot.version>4.1.1</spring-boot.version>` 属性
  （框架版本字面量全仓库唯一出现处，满足 FR-005），`dependencyManagement` 导入
  `spring-boot-dependencies`；未来第三方版本一律在此声明。
- `ecboot-parent`：`dependencyManagement` 导入 `ecboot-dependencies`，所有子模块经
  继承获得 BOM 管理的版本。
- `start`：`<parent>` 从 spring-boot-starter-parent 改为 ecboot-parent（直接继承父级，
  不再"直接继承上游框架父级"，符合 FR-005）；其现有 Boot starters 全部保留。

**Rationale**: 框架版本经 `ecboot-parent` 继承链获取（澄清 Q1 明确允许"继承链 /
导入清单"两种途径）；`spring-boot-starter-parent` 的 pluginManagement（spring-boot
编译/打包/原生插件版本）与编译默认值得以完整保留，`start` 零额外配置；BOM 导入
位于 ecboot-parent 自身的 dependencyManagement，按 Maven "就近优先"规则压过祖父级
继承的清单——`dependencies` 中的声明对第三方库具有最终仲裁权（满足 US2 场景 2）。

**Alternatives considered**:
- ecboot-parent 脱离 starter-parent、自管全部插件版本 → 插件版本与 BOM 的
  spring-boot.version 属性跨文件重复（导入不传播属性），制造双版本源，违反 FR-005 精神。
- 各模块各自导入 ecboot-dependencies BOM（不经父级）→ 10 处重复声明，违反 DRY 与宪法 V。
- start 保留 starter-parent 直继承（澄清 Q1 选项 B）→ 已被用户否决。

> **实现修订（2026-09-17）**：BOM 模块不继承 ecboot-parent（独立声明 GA），否则
> "根导入 BOM + BOM 继承根"构成 scope=import 自环。Maven 3.9.16 已支持 reactor
> 内 import 解析（无需预安装 BOM）。模块 groupId 归一为 `org.juling.ecboot`
> （与根 POM 一致）。

## D3: 依赖方向强制机制（FR-008、澄清 Q2 内置强制）

**Decision**: 采用 `maven-enforcer-plugin` 的 `bannedDependencies` 规则，按模块
白名单实现分层方向（允许矩阵见 data-model.md）。配置集中放置在根 POM 的
`pluginManagement`，规则内容通过每模块自定义属性（如 `enforcer.banned.excludes`）
参数化——各业务模块只需一行属性声明自己的禁止清单，执行绑定统一挂在根级
`validate` 阶段（构建生命周期最早阶段，满足"任何一次构建即校验"）。

**Rationale**: enforcer 是 Maven 生态事实标准的构建守护插件，纯声明式、无脚本；
`validate` 阶段在编译前失败，fail-fast 且错误信息含违规 GA 坐标；循环依赖另由
reactor 原生检测兜底（"cyclic reference" 构建失败，FR-007 双保险）。

**Alternatives considered**:
- 独立校验脚本（scripts/ 下 exec 绑定）→ 需维护自研脚本 + exec 插件，复杂度更高。
- 仅 CI 检查 → 仓库无 CI，且违反澄清 Q2 的"内置强制"。
- maven-archetype/Gradle 重构 → 推倒重来，违反宪法 II/V。

## D4: 版本纪律的机械强制边界（FR-005、US2 场景 3）

**Decision**: 三层机制——
1. **缺版本即失败（天然）**: 全部第三方版本由 BOM 管理，模块新增未登记依赖且不写
   版本号 → Maven 在构建期直接报错 "`dependencies.dependency.version` ... is
   missing"，错误即指向"去 dependencies 声明"（满足 US2 场景 3 主路径）。
2. **enforcer 附加规则**: `banDuplicatePomDependencyVersions`（同 POM 重复声明失败）+
   `reactorModuleConvergence`（全 reactor 模块版本必须一致）。
3. **显式硬编码版本的检测**: 接受为机制的已知边界——enforcer 原生规则无法直接禁止
   `<version>` 显式声明。处置：BOM 白名单 + 评审约束 + `banDuplicatePomDependencyVersions`；
   **不**引入第三方 enforcer 扩展（如 pedantic-pom-enforcer）作机械检测。

**Rationale**: vanilla enforcer 覆盖了方向违规与未登记版本两大主路径（FR-007/008、
US2.3 可全机械化验证）；硬编码版本属低频越轨行为，引入冷门第三方构建插件的
供应链成本高于收益（宪法 V）。此边界已在 quickstart 的验证场景中如实反映。

**Alternatives considered**:
- pedantic-pom-enforcer（可机械禁止显式版本）→ 第三方构建插件、维护不活跃，风险>收益。
- 自研 enforcer 规则 → 维护成本最高，违反宪法 V。

## D5: start 的依赖与构建配置迁移

**Decision**: `start/pom.xml` 仅改 `<parent>` 指向 ecboot-parent，**现有依赖清单
（12 个 starter + flyway-mysql、mysql-connector-j、devtools、docker-compose、lombok
及 12 个 test starter）全部保留**——Boot starters 是部署物的运行时组成，非业务模块，
不违反 FR-001（其约束对象是 org.juling 业务构件）。原 POM 中的编译插件配置
（lombok/ConfigurationProcessor 注解处理器路径）、hibernate 增强执行、testCompile
处理器路径**上移至根 POM** 供全模块继承；spring-boot-maven-plugin 在 start 保留
声明（repackage 执行由 starter-parent 的 pluginManagement 提供）。

**Rationale**: 上移后所有模块统一获得 Lombok 处理与 Hibernate 增强，避免后续业务
模块重复配置；start 自身 POM 瘦身为"依赖清单 + 应用插件"，符合装配模块定位。

**Alternatives considered**:
- 注解处理器/增强配置留在 start → 其他模块引入 Lombok 时需各自复制配置。
- 删除 start 的 starters 移到 services 层 → 改变运行时装配职责，超出本特性范围。

## D6: 插件版本集中（澄清 Q3 拆分治理）

**Decision**: 根 POM `pluginManagement` 新增/覆盖的插件版本仅 **maven-enforcer-plugin**
一项（选用当时稳定版 3.x）；其余插件（compiler、spring-boot、hibernate-enhance、
native、surefire 等）版本全部由继承链上的 spring-boot-starter-parent 托管，零声明。
业务模块 POM 不出现任何 `<version>`（依赖或插件）。

**Rationale**: 框架插件与框架版本同源（starter-parent），天然一致；新增插件唯一
需要显式版本，且声明点唯一（根父 POM），满足"插件版本集中在根父 POM"。

**Alternatives considered**:
- 插件版本也进 dependencies BOM → 拆分治理被打破（澄清 Q3 已否决）。
- 插件版本就地声明 → 版本散落，违反 FR-005。

## 未决项

无——所有 Phase 0 未知项已解决。
