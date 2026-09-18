# Implementation Plan: 社交电商能力扩展——数据库 Schema 设计（P0/P1/P2）

**Branch**: `002-social-commerce-expansion` | **Date**: 2026-09-18 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-social-commerce-expansion/spec.md`

## Summary

为 P0（电商闭环）、P1（社交差异化）、P2（增长运营）共 12 个域交付数据库 Schema 扩展：以 Flyway 增量迁移（V11~V21，按域一文件）新增约 28 张表、改造 `user` 表（手机号加密）并为 `trade_order`/`trade_order_item` 增列承载拼团/积分/满减；同步更新设计文档（`docs/schema-design.md`）、领域术语表（`CONTEXT.md`）并新增 ADR-0003（分销关系链两级封顶）。**不含任何应用代码**。

## Technical Context

**Language/Version**: SQL（MySQL 8.4 LTS 方言）；文档 Markdown（中文，宪法 III）

**Primary Dependencies**: Flyway（`spring-boot-starter-flyway`，随 `ecboot-start` 启动自动执行）；既有 V1~V10 迁移与 `docs/schema-design.md` 全局约定（金额/幂等/软删/命名）

**Storage**: MySQL 8.4 LTS（InnoDB、utf8mb4 / utf8mb4_0900_ai_ci），本地由 `apps/api/compose.yaml` 提供

**Testing**: 结构验证三板斧——① 迁移可执行（在已应用 V1~V10 的库上顺序执行 V11+，零报错）；② 约束注入（违反唯一约束的插入必须被数据库拒绝，以报错为证）；③ 校验 SQL（优惠分摊恒等式等，见 quickstart.md）

**Target Platform**: `apps/api/ecboot-start/src/main/resources/db/migration/`（Flyway 默认位置，宪法工程约束）

**Project Type**: 数据库 Schema 迁移集（纯 DDL/迁移脚本 + 设计文档）

**Performance Goals**: N/A（结构设计任务；仅保证高频查询路径索引齐备，与既有索引规范一致）

**Constraints**: 只允许 V11+ 增量（唯一既有表改造 = `user` 手机号加密，FR-003/FR-004）；金额一律 `DECIMAL(10,2)`、账户余额列不设无符号（可负）、枚举 TINYINT+注释、无物理外键、主数据软删/交易单据与日志永不删除

**Scale/Scope**: 新增 28 张表（P0:7、P1:13、P2:8），改造 1 张（`user`），为 2 张既有表增列（`trade_order`、`trade_order_item`）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束 | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ 通过 | 不新增顶层模块、不引入分布式；迁移文件落位 `ecboot-start` 资源目录，分层不变 |
| II 统一技术栈 | ✅ 通过 | 零新依赖——纯 SQL 迁移 + 既有 Flyway |
| III 中文优先 | ✅ 通过 | 迁移注释、设计文档、术语表全部中文；表名/列名（代码标识符）英文 |
| IV 可验证交付 | ✅ 通过 | 验证命令与预期输出落在 quickstart.md（迁移执行/约束注入/校验 SQL），以输出为证 |
| V 简单优先 | ✅ 通过 | 不建新模块/抽象层；全部落位既有目录与既有约定 |
| 工程约束（1.1.0） | ✅ 通过 | 迁移位于 `db/migration` 默认位置；不改 POM；不动包管理器；`user` 改造为增量 ALTER |

**Phase 1 复查**：✅ 设计产物（data-model/contracts/quickstart）未引入新模块、新依赖、新治理规则；全部通过，无违例需记录。

## Project Structure

### Documentation (this feature)

```text
specs/002-social-commerce-expansion/
├── plan.md              # 本文件 (/speckit-plan)
├── research.md          # Phase 0 研究与决策
├── data-model.md        # Phase 1 实体/字段/关系/状态机
├── quickstart.md        # Phase 1 验证指南
├── contracts/           # Phase 1 契约（迁移清单与约束契约）
│   └── schema-contracts.md
└── tasks.md             # /speckit-tasks 生成（本命令不创建）
```

### Source Code (repository root)

```text
apps/api/ecboot-start/src/main/resources/db/migration/   # Flyway 默认位置
├── V1__user_domain.sql … V10__admin_domain.sql          # 既有（26 表，不可改动）
├── V11__privacy_user_security.sql    # P0 user 改造(手机号加密) + user_login_log
├── V12__product_review.sql           # P0 product_review
├── V13__favorite_footprint.sql       # P0 user_favorite + user_footprint
├── V14__notify.sql                   # P0 notify_task + user_message
├── V15__logistics_company.sql        # P0 logistics_company
├── V16__distribution.sql             # P1 分销 8 表（关系链/推广员/规则/记录/账户/流水/提现/邀请）
├── V17__group_buy.sql                # P1 拼团 3 表 + trade_order 增列
├── V18__flash_sale.sql               # P1 秒杀 2 表
├── V19__point_level.sql              # P2 积分 2 表 + 等级 1 表 + trade_order/item 积分列
├── V20__full_reduction.sql           # P2 满减 3 表 + trade_order/item 满减列
└── V21__risk_control.sql             # P2 风控 2 表

docs/
├── schema-design.md                  # 按域增补（文件清单/决策/索引）
└── adr/0003-distribution-two-level-cap.md   # 新 ADR：分销关系链两级封顶

CONTEXT.md                            # 术语表增补（评价/收藏/足迹/通知/物流/关系链/佣金/账户/提现/拼团/秒杀/积分/等级/满减/风控）
```

**Structure Decision**: 迁移按"一域一文件、P0→P1→P2 顺序编号"（依据 research.md D1）；`trade_order`/`trade_order_item` 的增列分散到所属域的迁移文件内（V17/V19/V20），保证按域独立评审与回滚。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违例。唯一接近边界的设计是"满减范围用关系表（3 张表）而非 JSON 列"——不属于宪法违例，属 FR-020 的"高效命中查询"要求，已在 research.md D3 记录取舍。
