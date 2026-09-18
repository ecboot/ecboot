---
description: "Task list for 社交电商能力扩展——数据库 Schema 设计"
---

# Tasks: 社交电商能力扩展——数据库 Schema 设计

**Input**: Design documents from `/specs/002-social-commerce-expansion/`（plan.md、spec.md、research.md、data-model.md、contracts/schema-contracts.md、quickstart.md）

**Prerequisites**: plan.md (required), spec.md (required)

**Tests**: spec 未要求测试代码任务；验证以 quickstart.md 五场景的可执行 SQL 断言代替（每故事一个验证任务，输出留证，宪法 IV"以输出为证"）。

**Organization**: 按用户故事分组（US1~US11 ↔ 迁移 V11~V21 一一对应），每故事两任务（编写→执行验证），独立可交付。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件、无依赖）。**注意**：所有"执行迁移"任务共享同一 MySQL 实例（Flyway 锁），一律不标 [P]、串行执行。
- **[Story]**: 归属用户故事（US1~US11）
- 迁移文件均位于 `apps/api/ecboot-start/src/main/resources/db/migration/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 迁移基线环境就绪

- [ ] T001 启动本地 MySQL（`cd apps/api && docker compose up -d mysql`）并应用基线迁移 V1~V10（启动 `./mvnw spring-boot:run -pl ecboot-start -am` 或 mysql 客户端手工按序执行），核对 `flyway_schema_history` 10 行 success=1、业务表 26 张
- [ ] T002 [P] 核对迁移编号基线：确认 `db/migration/` 下最高版本为 V10、V11~V21 未被占用；将 data-model.md 各域字段清单标记为编写依据（不改动设计文档）

**Checkpoint**: 基线库就绪，编号无冲突

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 各域文档同步的共享锚点（本特性无代码基建，仅此一项跨故事前置）

- [ ] T003 在 `docs/schema-design.md` 文件清单表后新增"增量迁移总览（V11~V21）"小节：11 个文件的版本号/域/表名清单与执行前提（照抄 contracts/schema-contracts.md §1），作为后续各域决策增补的挂载点

**Checkpoint**: 文档锚点就绪，用户故事可开始

---

## Phase 3: User Story 1 - 隐私与账号安全数据模型 (Priority: P1) 🎯 MVP

**Goal**: `user` 手机号密文化（phone_hash 唯一检索）+ 用户登录日志表（V11）

**Independent Test**: quickstart 场景四——`uk_phone_hash` 存在、`uk_phone` 已删、`phone_hash CHAR(64) NOT NULL`、无明文片段列

### Implementation for User Story 1

- [ ] T004 [US1] 编写 `apps/api/ecboot-start/src/main/resources/db/migration/V11__privacy_user_security.sql`：`ALTER TABLE user`（`phone` 列注释改"手机号密文"并放宽至 VARCHAR(256)、`ADD phone_hash CHAR(64) NOT NULL`、`DROP INDEX uk_phone`、`ADD UNIQUE KEY uk_phone_hash(phone_hash)`）+ 新建 `user_login_log`（字段按 data-model P0/V11，`idx(user_id, created_at)`）；文件头注释含"存量数据环境须先执行应用侧密文化任务"警示（research D2）
- [ ] T005 [US1] 执行 V11 并运行 quickstart 场景四断言（SHOW INDEX / SHOW COLUMNS），输出留证到 `specs/002-social-commerce-expansion/` 验证记录（可追加至本文件 Notes 或单独 baseline 文件）

**Checkpoint**: US1 独立可验证——合规红线第 4 条的 Schema 侧就位

---

## Phase 4: User Story 2 - 商品评价数据模型 (Priority: P1)

**Goal**: 订单项粒度评价表，一项一评（V12）

**Independent Test**: quickstart 场景二②——同 order_item_id 二次插入报 ERROR 1062

- [ ] T006 [P] [US2] 编写 `.../db/migration/V12__product_review.sql`：`product_review` 全字段（含快照列 `spu_name`/`sku_specs`，research D9；追评/回复各一次的承载列；`audit_status` 状态机注释）、`UNIQUE(order_item_id)`、`idx(spu_id, audit_status, created_at)`
- [ ] T007 [US2] 执行 V12 并跑场景二②注入验证（1062 留证）

---

## Phase 5: User Story 3 - 收藏与足迹数据模型 (Priority: P2)

**Goal**: 用户×商品唯一关系的收藏/足迹两表（V13）

**Independent Test**: quickstart 场景二①——重复收藏报 1062；足迹表无 deleted、含 idx(last_view_at)

- [ ] T008 [P] [US3] 编写 `.../db/migration/V13__favorite_footprint.sql`：`user_favorite`（UNIQUE(user_id, spu_id)、软删）+ `user_footprint`（UNIQUE(user_id, spu_id)、`last_view_at`、无 deleted、`idx(last_view_at)`）
- [ ] T009 [US3] 执行 V13 并跑场景二①验证（1062 留证）

---

## Phase 6: User Story 4 - 消息通知数据模型 (Priority: P2)

**Goal**: 一渠道一行的异步通知任务表 + 站内消息表（V14，research D4）

**Independent Test**: DESC 两表核对状态机枚举注释（10/20/30/40/50）与重试字段齐备

- [ ] T010 [P] [US4] 编写 `.../db/migration/V14__notify.sql`：`notify_task`（channel/biz_type/template_code/params JSON/status 状态机/retry_count/next_retry_time，`idx(status, next_retry_time)`、`idx(user_id, created_at)`、`idx(biz_no)`）+ `user_message`（title/content/is_read，`idx(user_id, is_read, created_at)`）
- [ ] T011 [US4] 执行 V14，`SHOW CREATE TABLE` 核对两表状态机注释与索引，留证

---

## Phase 7: User Story 5 - 物流公司字典 (Priority: P2)

**Goal**: 承运方字典表（V15）

**Independent Test**: 插入一启用一停用样例行，字典可用且停用保留

- [ ] T012 [P] [US5] 编写 `.../db/migration/V15__logistics_company.sql`：`logistics_company`（`code VARCHAR(32) UNIQUE`、name、status、tracking_rule、软删）
- [ ] T013 [US5] 执行 V15 并插入两条样例（启用/停用各一）验证，留证

---

## Phase 8: User Story 6 - 分销推广数据模型 (Priority: P3)

**Goal**: 两级封顶关系链 + 佣金/账户/提现 8 表 + ADR-0003（V16）

**Independent Test**: quickstart 场景二⑤（`user_relation` 仅直接上级单列，结构上无法表达三级）+ 场景五 R2（`user_account.balance` 有符号）

- [ ] T014 [US6] 编写 `.../db/migration/V16__distribution.sql`：8 表按 data-model P1/V16（`user_relation` 仅 `inviter_id` 单列+UNIQUE(user_id)；`commission_rule` UNIQUE(scope_type, scope_id)；`commission_record` 含冲销关联列与四态状态机；`user_account` 可负 balance+frozen；`account_log` 只追加含 balance_after；`withdraw_order` UNIQUE(withdraw_channel, channel_order_no)+六态状态机；`invite_record` UNIQUE(new_user_id)；`distribution_user`）
- [ ] T015 [P] [US6] 编写 `docs/adr/0003-distribution-two-level-cap.md`：两级封顶决策（背景=禁止传销条例红线、决策=仅直接上级单列、被拒方案=多级路径表、后果=扩展层级须修宪级别评审），引用 spec SC-002 与 contracts 契约 R1
- [ ] T016 [US6] 执行 V16 并跑场景二⑤结构审查 + 场景五 R2 可负断言，留证

---

## Phase 9: User Story 7 - 拼团数据模型 (Priority: P3)

**Goal**: 拼团活动/团/团员 3 表 + 订单关联列（V17，research D10）

**Independent Test**: quickstart 场景二③（同团同用户 1062）+ 普通订单零影响（可空列、既有索引未变）

- [ ] T017 [US7] 编写 `.../db/migration/V17__group_buy.sql`：`group_buy_activity`（group_price/group_size/时段/per_limit）、`group_buy_team`（状态机 1/2/3、`idx(status, expire_time)`）、`group_buy_team_member`（UNIQUE(team_id, user_id) + UNIQUE(order_no)）、`ALTER TABLE trade_order ADD group_buy_team_id BIGINT NULL + KEY idx_group_team`
- [ ] T018 [US7] 执行 V17 并跑场景二③ + 零影响断言（`SHOW COLUMNS trade_order LIKE 'group_buy_team_id'` 可空），留证

---

## Phase 10: User Story 8 - 秒杀数据模型 (Priority: P3)

**Goal**: 秒杀活动/场次商品 2 表，活动库存分账（V18，research D7）

**Independent Test**: quickstart 场景五 R4——stock_count/sold_count 在 flash_sale_item，与 inventory 零耦合

- [ ] T019 [P] [US8] 编写 `.../db/migration/V18__flash_sale.sql`：`flash_sale_activity`（时段/状态）+ `flash_sale_item`（flash_price DECIMAL(10,2)、stock_count/sold_count、per_limit、UNIQUE(activity_id, sku_id)、`idx(sku_id)`）；文件尾注释活动库存条件更新语义（同构 ADR-0001）
- [ ] T020 [US8] 执行 V18 并跑场景五 R4 分账断言，留证

---

## Phase 11: User Story 9 - 积分与会员等级数据模型 (Priority: P4)

**Goal**: 积分双账本（可负账户+流水）+ 等级规则 + user/订单积分列（V19，research D5/D8）

**Independent Test**: `point_account.balance` 有符号；`user.level` 可空；trade_order 新增 point_amount/point_used 带默认值

- [ ] T021 [US9] 编写 `.../db/migration/V19__point_level.sql`：`point_account`（balance INT 有符号、UNIQUE(user_id)）+ `point_log`（biz_type 五态、balance_after、idx(user_id, id)）+ `user_level_rule`（UNIQUE(growth_threshold)、benefits JSON）+ `ALTER user ADD growth_value INT UNSIGNED DEFAULT 0 / level TINYINT NULL` + `ALTER trade_order ADD point_amount DECIMAL(10,2) NOT NULL DEFAULT 0 / point_used INT UNSIGNED NOT NULL DEFAULT 0` + `ALTER trade_order_item ADD point_amount DECIMAL(10,2) NOT NULL DEFAULT 0`
- [ ] T022 [US9] 执行 V19，DESC 断言可负/可空/默认值三类列形态，留证

---

## Phase 12: User Story 10 - 满减活动数据模型 (Priority: P4)

**Goal**: 满减活动/档位/范围 3 表 + 订单优惠三构成列（V20，research D3/D8）

**Independent Test**: quickstart 场景二④（同活动同门槛 1062）+ 场景三恒等式校验（需 V19 已执行）

- [ ] T023 [US10] 编写 `.../db/migration/V20__full_reduction.sql`：`promotion_activity`、`promotion_activity_ladder`（UNIQUE(activity_id, threshold_amount)）、`promotion_activity_scope`（scope_type 1全场/2分类/3商品、target_id NULL=全场、UNIQUE 三元组、`idx(scope_type, target_id)`）+ `ALTER trade_order ADD coupon_amount / full_reduction_amount DECIMAL(10,2) NOT NULL DEFAULT 0 / promotion_activity_id BIGINT NULL` + `ALTER trade_order_item ADD coupon_amount / full_reduction_amount`
- [ ] T024 [US10] 执行 V20，跑场景二④；随后按 quickstart 场景三手工 INSERT 样例订单并执行恒等式校验 SQL（依赖 T021 的 point 列），两断言留证

---

## Phase 13: User Story 11 - 风控数据模型 (Priority: P4)

**Goal**: 风控规则/事件 2 表（V21）。**站内搜索零建表**为 spec 显式决定——本故事不含搜索任务。

**Independent Test**: `risk_record` 含 object_type 多态关联与 appeal_status 申诉四态

- [ ] T025 [P] [US11] 编写 `.../db/migration/V21__risk_control.sql`：`risk_rule`（rule_type 四类、condition_expr、action、软删）+ `risk_record`（只追加；rule_id、object_type/object_no、action、appeal_status 0-3、`idx(user_id, created_at)`、`idx(object_no)`）
- [ ] T026 [US11] 执行 V21，SHOW CREATE TABLE 核对申诉枚举与双索引，留证

---

## Phase 14: Polish & Cross-Cutting Concerns

- [ ] T027 完成文档同步三件套（契约 §5 逐项勾验）：`docs/schema-design.md` 按域增补决策/索引说明、文件清单扩至 V21、表总数 26→54；`CONTEXT.md` 增补 15 个术语（清单见契约 §5）；确认 `docs/adr/0003` 已在 ADR 索引可发现
- [ ] T028 全量回归：全新空库重放 V1~V21，跑 quickstart 场景一（flyway 21 行 success=1、业务表 54）与场景三恒等式终验，输出留证
- [ ] T029 收尾：`git status` 全量核对交付物（11 迁移 + 3 文档 + tasks/spec 系列），验证输出归档至特性目录，按宪法准备中文 Conventional Commits 分批提交（迁移与文档分批，如 `feat: 新增P0闭环域数据库迁移V11-V15`）

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → **用户故事（Phase 3~13）**：故事间无代码依赖，仅共享同一迁移库——**执行/验证任务必须串行**（Flyway 锁），编写任务可并行
- **Polish (Phase 14)**：T027 依赖全部编写任务；T028 依赖全部执行任务；T029 最后

### User Story Dependencies

- US1~US11 相互独立，可任意顺序交付；例外：**T024 的恒等式校验依赖 T021（V19 point 列）已执行**——若 US10 先于 US9，恒等式校验延后到 T028 全量回归补跑
- 优先顺序按 spec：P1（US1、US2）→ P2（US3、US4、US5）→ P3（US6、US7、US8）→ P4（US9、US10、US11）

### Parallel Opportunities

- 编写类任务（T004/T006/T008/T010/T012/T015/T019/T025，不同文件）可全部并行
- 同故事内"编写→执行验证"严格串行；跨故事执行验证串行（共享库）

## Parallel Example: 编写阶段

```bash
# 可同时派发（互不相同文件、无依赖）：
Task: "编写 V12__product_review.sql"      # US2
Task: "编写 V13__favorite_footprint.sql"  # US3
Task: "编写 V14__notify.sql"              # US4
Task: "编写 V15__logistics_company.sql"   # US5
# 执行验证（T007/T009/...）待各自编写完成后串行入库
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1（基线库）→ Phase 2（文档锚点）→ Phase 3（US1：手机号加密 + 登录日志）
2. **STOP and VALIDATE**：场景四断言通过即可交付——合规红线项最先落地

### Incremental Delivery

每完成一个故事 = 一个独立可评审/可回滚的迁移文件 + 验证留证 + 一次提交；P0（V11~V15）→ P1（V16~V18）→ P2（V19~V21）三批可各自作为上线单元。

---

## Notes

- 每任务或逻辑任务组完成后提交一次（中文 Conventional Commits，宪法 III），如 `feat: 新增分销域数据库迁移V16与两级封顶ADR`
- 验证一律以命令输出为证（宪法 IV），留证文件放 `specs/002-social-commerce-expansion/`
- 迁移文件不改动 V1~V10；既有列语义不变（契约 §1 兼容性）
- 本 tasks.md 交由 `/speckit-superpowers-bridge` 接管执行时，遵守 handoff 约定（guard 阻止 `speckit.implement` 期间改契约产物）
