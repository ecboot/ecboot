# Research: 社交电商能力扩展——Schema 设计（Phase 0）

2026-09-18。spec 无 [NEEDS CLARIFICATION] 遗留；本文件收敛 plan 阶段的 12 项设计决策（含 clarify 阶段 Deferred 的 3 项）。

## D1. 迁移文件划分：一域一文件，V11~V21 顺序编号

- **Decision**: 11 个迁移文件按 P0→P1→P2 编号（见 plan.md 结构树）；`trade_order`/`trade_order_item` 增列随所属域文件。
- **Rationale**: spec 假设"按域独立成文件、可按域评审与回滚"；域是评审与回滚的自然单元。
- **Alternatives**: 一层一文件（3 个大文件，评审粒度过粗、回滚连带无辜域）；一表一文件（28 个文件，编号噪音）。

## D2. user 手机号加密：改造 `phone` 列语义 + 新增 `phone_hash`，单次 ALTER

- **Decision**: V11 内 `ALTER TABLE user`：`phone` 列注释与语义改为"手机号密文（加密存储）"，新增 `phone_hash CHAR(64) NOT NULL`（加盐 SHA-256），`DROP INDEX uk_phone`，`ADD UNIQUE KEY uk_phone_hash (phone_hash)`。不新增并列密文列、不留明文片段列（spec FR-004：仅精确检索）。
- **Rationale**: 项目脚手架期、库中无用户数据（V1 尚未在任何真实环境执行），单次结构改造即可，无双轨期成本（spec 假设"不提供长期双轨"的依据即此）。
- **Alternatives**: 新增 `phone_cipher` 列与旧列共存再删（为有存量数据的环境设计，当前无数据属过度设计）；可逆加密（spec 澄清已拒绝）。
- **注**: 加密算法/盐策略属实现任务；若未来在有存量数据的环境执行，须先跑应用侧密文化任务再执行本迁移——写入迁移文件头注释。

## D3. 满减范围：关系表 `promotion_activity_scope`（非 JSON）

- **Decision**: `promotion_activity_scope(activity_id, scope_type TINYINT 1全场/2分类/3商品, target_id BIGINT NULL)`，`UNIQUE(activity_id, scope_type, target_id)`，`idx(scope_type, target_id)` 支持反查命中。
- **Rationale**: FR-020 要求高效命中查询（"哪些活动覆盖此商品/分类"是下单链路热路径）；JSON 无法建此类索引。全场用 scope_type=1 + target_id=NULL 单行表达。
- **Alternatives**: JSON 列存范围集合（无法索引、命中全表扫）；每活动多张范围表（分类表+商品表，表数翻倍无增益）。

## D4. 通知任务粒度：一渠道一行

- **Decision**: `notify_task` 每渠道一行（user_id, channel TINYINT 1小程序订阅消息/2短信/3站内信, template_code, params JSON, status, retry_count, next_retry_time）；渠道被用户关闭则不建行（站内信例外，必达）。
- **Rationale**: 重试计数、渠道状态、用户订阅偏好都是渠道粒度；一行多渠道（JSON 渠道数组）会让重试语义纠缠。
- **Alternatives**: 一行多渠道（rejected 如上）；每模板一行（与渠道粒度冲突）。

## D5. 成长值承载：`user` 表列（非独立账本）

- **Decision**: `user.growth_value INT UNSIGNED NOT NULL DEFAULT 0` 累计列 + `user.level TINYINT NULL`（可空，未启等级体系零影响，FR-019）；成长值只增不减（退款不扣），故无流水对账诉求，等级变更历史由后台操作审计承载。
- **Rationale**: clarify 结论"双账本分离"指积分（可消耗，有账户+流水）与成长值（只增不减）分离；成长值自身简单累计，独立账本表属过度设计。
- **Alternatives**: 独立 growth 账本+流水表（成长值无消耗无回退，流水无对账价值）；由积分流水 SUM 派生（clarify 已拒绝，聚合查询随流水膨胀劣化）。

## D6. 佣金规则作用域：单表两作用域 + UK，命中"商品>分类"

- **Decision**: `commission_rule(scope_type TINYINT 1分类/2商品, scope_id BIGINT, level1_rate DECIMAL(5,2), level2_rate DECIMAL(5,2), status, deleted…)`，`UNIQUE(scope_type, scope_id)`；命中优先级应用层实现"商品>分类，未命中不计佣"（clarify Q1 结论），比例存百分数值（如 10.00 = 10%）。
- **Rationale**: 两作用域同构（类型+目标+两比例），单表最简；优先级是查询逻辑不是结构逻辑，无需冗余层级列。
- **Alternatives**: 分类表+商品表两张（同构拆表无增益）；三级作用域（clarify 已拒绝，过度设计）。

## D7. 秒杀活动库存分账：`flash_sale_item` 自带限量/已售

- **Decision**: `flash_sale_item(activity_id, sku_id, flash_price DECIMAL(10,2), stock_count INT UNSIGNED, sold_count INT UNSIGNED, per_limit INT UNSIGNED, UNIQUE(activity_id, sku_id))`；下单扣减/取消回补均作用本表（`UPDATE ... WHERE stock_count - sold_count >= n` 条件更新，与既有库存模型同构）；订单项快照 `price` 直接承载秒杀价（FR-016）。
- **Rationale**: 活动库存与普通库存生命周期不同（活动结束即失效），分账避免污染 `inventory`；条件更新复用 ADR-0001 验证过的防超卖模式。
- **Alternatives**: 复用 `inventory` 加 bucket 维度（改动既有热表、耦合活动生命周期）。

## D8. 优惠三构成在订单上的表达：合计列 + 三明细列

- **Decision**: `trade_order` 新增 `coupon_amount`/`full_reduction_amount`/`point_amount`（均 DECIMAL(10,2) NOT NULL DEFAULT 0）+ `point_used INT UNSIGNED`（消耗积分点数）+ `promotion_activity_id BIGINT NULL`（命中满减活动）；既有 `promotion_amount` 保留且恒等于三明细之和；`trade_order_item` 对应新增三明细分摊列。恒等式校验 SQL 落 quickstart（FR-018）。
- **Rationale**: 既有 `promotion_amount`（合计）已被 V6 语义占用且向后兼容；明细列让退款回退（券退回/积分回退/满减重算）有独立落点，避免从 `user_coupon_id` 反推。
- **Alternatives**: 只存明细删合计列（破坏 V6 兼容）；只存合计（回退语义不可分）。

## D9. 评价商品快照：冗余商品标题列

- **Decision**: `product_review` 冗余 `spu_name VARCHAR(128)`、`sku_specs JSON` 快照（同 ADR-0002 快照原则），SPU 软删后评价自洽展示（spec 边界场景）。
- **Rationale**: 评价属用户内容+消费决策依据，商品删除不随删；引用列（spu_id/sku_id）保留用于聚合统计。
- **Alternatives**: 仅存关联 id（软删后展示悬挂）；评价随商品隐藏（spec 边界场景已拒绝）。

## D10. 拼团订单关联：`trade_order.group_buy_team_id` 可空列

- **Decision**: 团员表 `group_buy_team_member(team_id, user_id, order_no)` 三重唯一各司其职：`UNIQUE(team_id, user_id)`（一团一单，FR-015）、`UNIQUE(order_no)`（一单至多一团）、订单侧 `group_buy_team_id BIGINT NULL + idx`（普通订单零影响）；成团价经订单项 `price` 快照承载，无需新金额列。
- **Rationale**: 普通订单路径零改动（可空列 + 不命中索引）；约束放团员表而非订单表，域内自洽。
- **Alternatives**: 独立拼团订单表（分裂订单主链路，报表/售后全要 UNION）。

## D11. 提现幂等：渠道单号唯一 + 状态机条件更新

- **Decision**: `withdraw_order` 含 `channel_order_no VARCHAR(64) NULL`（打款渠道单号，`UNIQUE(withdraw_channel, channel_order_no)`）与状态枚举 10待审核/20审核通过/30打款中/40成功/50审核拒绝/60打款失败已回退；打款回调以条件更新抢占（`WHERE status=30`），重复回调 affected=0 直接 ACK——完全复用支付回调幂等总则（FR-013）。
- **Rationale**: 与 pay_order 同构的资金单据幂等模式已被设计验证，不发明新机制。
- **Alternatives**: 防重表（多一跳，无增益）。

## D12. 新 ADR-0003：分销关系链两级封顶

- **Decision**: 新增 `docs/adr/0003-distribution-two-level-cap.md`，记录"关系链仅 `inviter_id` 单列、结构上无法表达三级"的法规红线决策（禁止传销条例；满足 ADR 三条件：难逆转/事后看意外/真实取舍——被拒方案是更灵活的多级路径表）。
- **Rationale**: spec SC-002 要求评审确认"结构上无法表达三级"；ADR 是该结论的正式载体，六个月后有人想加"祖父列"时会被它拦住。

## 关键事实核对（研究依据）

- 既有迁移 V1~V10 已存在且未在任何真实环境执行（脚手架期）→ D2 单次改造成立。
- `trade_order.promotion_amount` 语义为"优惠总额合计"（V6）→ D8 兼容性成立。
- ADR-0001 条件更新模式、ADR-0002 快照/不可删原则 → D7/D9/D11 复用成立。
- flyway-mysql 已在 `ecboot-start` 运行时依赖 → 迁移随启动执行，无需新构建接线。
