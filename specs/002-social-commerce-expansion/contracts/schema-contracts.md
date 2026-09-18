# Schema Contracts: 社交电商能力扩展

本特性对外暴露的"契约"是**迁移脚本集本身及其结构承诺**——它们是应用层（后续实现任务）与数据之间唯一的接口。完整字段见 [data-model.md](../data-model.md)，此处是可被执行与审查的契约清单。

## 1. 迁移契约（执行顺序与前提）

| 版本 | 文件 | 域 | 前提 |
|---|---|---|---|
| V11 | `V11__privacy_user_security.sql` | 隐私安全 | V1 已应用；**硬性前置：仅可在空 `user` 表或已完成应用侧密文化的库上执行**（runbook 见 `docs/schema-design.md` §V11 存量环境执行顺序）；前置未满足时迁移主动失败（防呆，不污染数据） |
| V12 | `V12__product_review.sql` | 评价 | V6 |
| V13 | `V13__favorite_footprint.sql` | 收藏足迹 | V2 |
| V14 | `V14__notify.sql` | 通知 | V1 |
| V15 | `V15__logistics_company.sql` | 物流 | — |
| V16 | `V16__distribution.sql` | 分销（8 表） | V2/V6 |
| V17 | `V17__group_buy.sql` | 拼团（3 表+订单列） | V6 |
| V18 | `V18__flash_sale.sql` | 秒杀（2 表） | V2 |
| V19 | `V19__point_level.sql` | 积分等级（3 表+user/订单列） | V1/V6 |
| V20 | `V20__full_reduction.sql` | 满减（3 表+订单列） | V2/V6/V9 |
| V21 | `V21__risk_control.sql` | 风控（2 表） | V1 |
| V22 | `V22__review_fixes.sql` | 评审修复 | V20；满减全场行唯一性加固 + 关联列索引补齐 |

- 执行方式：随 `ecboot-start` 启动由 Flyway 自动执行（默认位置，零配置）。
- 兼容性：V1~V10 文件不改动；既有列语义不变；`trade_order`/`trade_order_item`/`user` 仅**可空增列或带默认值增列**，既有读写路径零感知。
- 回滚边界：按域一文件，回滚以域为单位（Flyway 不做 down，回滚=恢复备份+重放，见 quickstart 场景五附注）。

## 2. 唯一性契约（数据库层强制，应用层不可绕过）

| 约束 | 表 |
|---|---|
| 一项一评 | `product_review UNIQUE(order_item_id)` |
| 一用户一收藏 / 一足迹 | `user_favorite` / `user_footprint` `UNIQUE(user_id, spu_id)` |
| 手机号唯一（哈希） | `user UNIQUE(phone_hash)` |
| 一人一条关系链 | `user_relation UNIQUE(user_id)` |
| 佣金规则作用域唯一 | `commission_rule UNIQUE(scope_type, scope_id)` |
| 提现渠道单号幂等 | `withdraw_order UNIQUE(withdraw_channel, channel_order_no)` |
| 一新用户仅一次邀请激励 | `invite_record UNIQUE(new_user_id)` |
| 一团一单 / 一单一团 | `group_buy_team_member UNIQUE(team_id, user_id)`、`UNIQUE(order_no)` |
| 一活动一 SKU / 一门槛 | `flash_sale_item UNIQUE(activity_id, sku_id)`；`promotion_activity_ladder UNIQUE(activity_id, threshold_amount)` |
| 满减范围唯一 | `promotion_activity_scope UNIQUE(activity_id, scope_type, target_id_norm)`（V22：生成列 `target_id_norm=IFNULL(target_id,0)` 归一 NULL，全场单行由数据库强制） |
| 等级门槛唯一 | `user_level_rule UNIQUE(growth_threshold)` |

## 3. 结构性强制规则（评审断言，非索引可表达）

- **R1 分销两级封顶**（法规红线，ADR-0003）：`user_relation` 仅 `inviter_id` 单列，不存在祖父列/路径列/层级列——任何三级及以上关系在结构上无处落库（SC-002）。
- **R2 余额可负**：`user_account.balance`、`point_account.balance` 为**有符号**列（欠款抵扣语义），禁止 UNSIGNED。
- **R3 敏感字段最小化**：`user` 表无明文手机号列、无尾号片段列；密文列 + 哈希索引二列封顶。
- **R4 活动库存分账**：秒杀限量/已售只存在于 `flash_sale_item`，与 `inventory` 零耦合。
- **R5 佣金防重复计佣由应用层承担**（评审裁定）：冲销负记录需要同 `(order_item_id, beneficiary_user_id, level)` 键多行存在，数据库唯一键不可表达"原始记录唯一 + 冲销记录成对"；应用层在产生佣金时以条件插入（先查后插于同一事务）保证一原始一冲销。资金路径，禁止留白以外的静默变更。
- **R6 关系链环/自邀校验义务在应用层**：两级封顶（R1）是结构强制的"层级"约束；`user_id = inviter_id`（自邀）与 A↔B 互环在结构上可插入，MUST 由应用层校验拒绝。

## 4. 恒等式契约（校验 SQL 见 quickstart 场景三）

- 订单头：`promotion_amount ≡ coupon_amount + full_reduction_amount + point_amount`
- 订单项：各明细分摊列合计 ≡ 订单头对应明细列
- 支付等式（既有，回归）：`pay_amount ≡ total_amount - promotion_amount + freight_amount`
- **适用范围（V22 声明）**：恒等式仅约束 V19/V20 之后创建的订单；此前存量订单的三构成列为默认 0 而 `promotion_amount` 可能 >0，数据巡检须按迁移应用时间过滤增量。

## 5. 文档同步契约（交付物的一部分）

- `docs/schema-design.md`：文件清单扩至 V21、按域补充决策与索引说明、表总数 26→54。
- `CONTEXT.md`：增补 15 个术语（评价/收藏/足迹/通知任务/站内消息/物流公司/关系链/推广员/佣金规则/佣金记录/账户/提现单/拼团/秒杀/积分/成长值/满减/风控）。
- `docs/adr/0003-distribution-two-level-cap.md`：新 ADR（两级封顶）。
