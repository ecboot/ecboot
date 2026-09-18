# Data Model: 社交电商能力扩展（Phase 1）

2026-09-18。新增 28 表 + 改造 `user` + 增列 `trade_order`/`trade_order_item`。全局约定沿用 `docs/schema-design.md`（InnoDB、utf8mb4、前缀命名、`created_at`/`updated_at`、无物理外键、金额 DECIMAL(10,2)、TINYINT 枚举、主数据 `deleted` 软删/单据日志只追加）。下文仅列业务列与约束（时间戳三件套与常规含义不再重复）。

## P0 — 电商闭环

### V11 `user`（改造）+ `user_login_log`

**user** 改造（research D2）：
- `phone VARCHAR(256)` — 语义改为**手机号密文**（注释更新；宽度容纳密文）
- `phone_hash CHAR(64) NOT NULL` — 加盐 SHA-256（登录精确检索与唯一性）
- 索引：`DROP uk_phone`；`ADD UNIQUE uk_phone_hash(phone_hash)`
- 规则：无明文片段列；检索仅哈希精确命中（FR-004）

**user_login_log**（只追加）：`user_id BIGINT NOT NULL idx(user_id, created_at)`、`login_channel TINYINT 1小程序 2H5`、`login_status TINYINT 1成功 2失败`、`ip VARCHAR(45)`、`user_agent VARCHAR(512)`。

### V12 `product_review`

- `UNIQUE(order_item_id)` — 一项一评（FR-006）
- 列：`user_id`、`spu_id`/`sku_id`（idx 聚合统计）、`spu_name VARCHAR(128)`/`sku_specs JSON`（快照，D9）、`score TINYINT 1-5`、`content VARCHAR(1024)`、`images JSON`、`is_anonymous TINYINT`、`audit_status TINYINT 0待审 1通过 2驳回`、`reply_content VARCHAR(512)`/`reply_time DATETIME`（商家回复一次）、`extra_content VARCHAR(1024)`/`extra_time DATETIME`（追评一次）
- 索引：`idx(spu_id, audit_status, created_at)`（商品页评价列表）
- 状态机：`audit_status` 0→1/2；追评仅一次（应用层，字段唯一承载）

### V13 `user_favorite` / `user_footprint`

- **user_favorite**：`UNIQUE(user_id, spu_id)`（FR-009）；`deleted` 软删（收藏可恢复场景留软删）。
- **user_footprint**：`UNIQUE(user_id, spu_id)` + `last_view_at DATETIME NOT NULL`（重复浏览=更新时间，D 语义）；无 `deleted`、定期清理（90 天，应用任务）。物理删除即可 + `idx(last_view_at)` 供清理扫描。

### V14 `notify_task` / `user_message`

- **notify_task**（一渠道一行，D4）：`user_id idx(user_id, created_at)`、`channel TINYINT 1小程序订阅 2短信 3站内信`、`biz_type TINYINT 1订单 2营销 3售后`、`biz_no VARCHAR(32)`（关联业务单号 idx）、`template_code VARCHAR(64)`、`params JSON`、`status TINYINT 10待发送 20已发送 30失败 40达重试上限 50已跳过`、`retry_count INT UNSIGNED DEFAULT 0`、`next_retry_time DATETIME NULL`。
  - 状态机：10→20（成功终态）；10→30→(retry_count<上限)→10 / →40（终态）；渠道被关闭→50（跳过终态）。`idx(status, next_retry_time)` 供重试扫描。
- **user_message**（站内必达兜底）：`user_id idx(user_id, is_read, created_at)`、`title VARCHAR(128)`、`content VARCHAR(2048)`、`biz_type`/`biz_no`、`is_read TINYINT`。

### V15 `logistics_company`

- 主数据软删：`code VARCHAR(32) UNIQUE`、`name VARCHAR(64)`、`status TINYINT 1启用 0停用`、`tracking_rule VARCHAR(255)`（单号校验规则描述）。与 `trade_order.deliver_company`（存 code）衔接。

## P1 — 社交差异化

### V16 分销 8 表（research D6/D11，新 ADR-0003）

- **user_relation**：`user_id BIGINT UNIQUE`（一人一条关系，FR-010）、`inviter_id BIGINT NOT NULL idx`、`bind_channel TINYINT 1分享链接 2邀请码`、`bind_time DATETIME`。**仅直接上级单列——结构上无法表达三级**（法规红线）。
- **distribution_user**：`user_id UNIQUE`、`status TINYINT 1待审核 2通过 3冻结`、`apply_time`/`audit_time`、`deleted`。
- **commission_rule**：`scope_type TINYINT 1分类 2商品`、`scope_id BIGINT`、`level1_rate DECIMAL(5,2)`、`level2_rate DECIMAL(5,2)`（百分数，0-100）、`status`、`deleted`；`UNIQUE(scope_type, scope_id)`；命中优先级 商品>分类（应用层）。
- **commission_record**：`order_no VARCHAR(32) idx`、`order_item_id`、`beneficiary_user_id idx(user_id, status)`、`level TINYINT 1/2`、`base_amount DECIMAL(10,2)`（=订单项实付）、`rate DECIMAL(5,2)`、`amount DECIMAL(10,2)`、`status TINYINT 1待结算 2已结算 3已失效 4欠款冲销中`、`settle_time DATETIME`、`reversal_record_id BIGINT NULL`（冲销关联原记录）。
  - 状态机：1→2（确认收货+7 天保护期满）；1→3（保护期退款）；2→4（结算后退款，产生冲销负记录）。
- **user_account**：`user_id UNIQUE`、`balance DECIMAL(10,2) NOT NULL DEFAULT 0`（**可负**——欠款抵扣）、`frozen DECIMAL(10,2) NOT NULL DEFAULT 0`（提现冻结）。
- **account_log**（只追加）：`user_id idx(user_id, id)`、`biz_type TINYINT 1佣金入账 2提现冻结 3提现完成 4提现回退 5冲销`、`amount DECIMAL(10,2)`（有符号）、`balance_after DECIMAL(10,2)`、`biz_no VARCHAR(32)`（关联 commission_record/withdraw_order）。
- **withdraw_order**：`withdraw_no VARCHAR(32) UNIQUE`、`user_id idx`、`amount DECIMAL(10,2)`、`withdraw_channel TINYINT 1微信商家转账`、`channel_order_no VARCHAR(64) NULL`、`UNIQUE(withdraw_channel, channel_order_no)`（幂等键，D11）、`status TINYINT 10待审核 20审核通过 30打款中 40成功 50审核拒绝 60打款失败已回退`、`audit_time`/`pay_time`/`fail_reason`。
  - 状态机：10→20/50；20→30（冻结余额）；30→40（解冻扣减，条件更新 `WHERE status=30`）/60（解冻回退）。
- **invite_record**：`new_user_id BIGINT UNIQUE`（一人仅被激励一次，FR-014）、`inviter_id idx`、`reward_type TINYINT 1优惠券`、`reward_ref BIGINT`（如 user_coupon.id）、`status TINYINT 1已发放`。

### V17 拼团 3 表 + 订单增列（research D10）

- **group_buy_activity**：`spu_id idx`、`name VARCHAR(64)`、`group_price DECIMAL(10,2)`、`group_size INT UNSIGNED`（成团人数）、`valid_start_at/valid_end_at DATETIME`、`per_limit INT UNSIGNED`、`status TINYINT 1启用 0停用`、`deleted`。
- **group_buy_team**：`activity_id idx(activity_id, status)`、`leader_user_id`、`status TINYINT 1拼团中 2已成团 3已解散`、`expire_time DATETIME`（`idx(status, expire_time)` 超时扫描）。
- **group_buy_team_member**：`team_id`、`user_id`、`order_no`；`UNIQUE(team_id, user_id)`（一团一单）、`UNIQUE(order_no)`（一单一团）。
- **trade_order 增列**：`group_buy_team_id BIGINT NULL, KEY idx_group_team(group_buy_team_id)`——普通订单零影响。

### V18 秒杀 2 表（research D7）

- **flash_sale_activity**：`name VARCHAR(64)`、`start_time/end_time DATETIME`、`status TINYINT 1启用 0停用`、`deleted`。
- **flash_sale_item**：`activity_id idx(activity_id, sku_id)`、`sku_id idx(sku_id)`（反查命中）、`flash_price DECIMAL(10,2)`、`stock_count INT UNSIGNED`、`sold_count INT UNSIGNED DEFAULT 0`、`per_limit INT UNSIGNED`、`UNIQUE(activity_id, sku_id)`。
  - 规则：可售=`stock_count-sold_count` 推导；扣减/回补条件更新同构 ADR-0001；秒杀价经 `trade_order_item.price` 快照承载。

## P2 — 增长运营

### V19 积分 2 表 + 等级 1 表 + 订单增列（research D5/D8）

- **point_account**：`user_id UNIQUE`、`balance INT NOT NULL DEFAULT 0`（**可负**——回退欠款，FR-017）。
- **point_log**（只追加）：`user_id idx(user_id, id)`、`biz_type TINYINT 1签到 2消费获得 3下单消耗 4退款回退 5分享获得`、`points INT`（有符号）、`balance_after INT`、`order_no VARCHAR(32)`。
- **user_level_rule**：`name VARCHAR(32)`、`growth_threshold INT UNSIGNED UNIQUE`（门槛唯一递增，FR-019）、`benefits JSON`（权益配置）、`status`、`deleted`。
- **user 增列**：`growth_value INT UNSIGNED NOT NULL DEFAULT 0`（只增不减）、`level TINYINT NULL`（未启等级为 NULL）。
- **trade_order 增列**：`point_amount DECIMAL(10,2) NOT NULL DEFAULT 0`、`point_used INT UNSIGNED NOT NULL DEFAULT 0`。
- **trade_order_item 增列**：`point_amount DECIMAL(10,2) NOT NULL DEFAULT 0`（行分摊）。

### V20 满减 3 表 + 订单增列（research D3/D8）

- **promotion_activity**：`name VARCHAR(64)`、`scope 由 scope 表表达`、`start_time/end_time DATETIME`、`status TINYINT 1启用 0停用`、`deleted`。
- **promotion_activity_ladder**：`activity_id`、`threshold_amount DECIMAL(10,2)`、`discount_amount DECIMAL(10,2)`、`UNIQUE(activity_id, threshold_amount)`（一活动同门槛唯一，FR-020）。
- **promotion_activity_scope**：`activity_id`、`scope_type TINYINT 1全场 2分类 3商品`、`target_id BIGINT NULL`（全场=NULL）、`UNIQUE(activity_id, scope_type, target_id)`、`idx(scope_type, target_id)`（下单热路径反查命中）。
- **trade_order 增列**：`coupon_amount DECIMAL(10,2) NOT NULL DEFAULT 0`、`full_reduction_amount DECIMAL(10,2) NOT NULL DEFAULT 0`、`promotion_activity_id BIGINT NULL`（命中活动）。
- **trade_order_item 增列**：`coupon_amount`、`full_reduction_amount`（行分摊）。

**恒等式**（FR-018，校验 SQL 见 quickstart）：`promotion_amount = coupon_amount + full_reduction_amount + point_amount`（订单头与订单项两层各自成立，行分摊合计=头合计，尾差记末行沿用既有规则）。

### V21 风控 2 表

- **risk_rule**：`name VARCHAR(64)`、`rule_type TINYINT 1黑名单 2高频下单 3异常领券 4佣金套利特征`、`condition_expr VARCHAR(512)`（规则条件描述）、`action TINYINT 1拦截 2标记`、`status TINYINT 1启用 0停用`、`deleted`。
- **risk_record**（只追加）：`user_id idx(user_id, created_at)`、`rule_id`、`object_type TINYINT 1订单 2优惠券 3提现 4售后`、`object_no VARCHAR(32)` idx、`action TINYINT 1拦截 2标记`、`appeal_status TINYINT 0无 1申诉中 2申诉通过 3申诉驳回`、`remark VARCHAR(255)`（FR-021）。

## 关系总览（增量部分）

```text
user 1─1 user_relation ──N user (inviter)
user 1─N user_login_log / user_favorite / user_footprint / product_review /
        notify_task / user_message / distribution_user / invite_record /
        point_account / point_log / user_account / account_log / withdraw_order / risk_record
product_spu 1─N product_review / user_favorite / user_footprint /
             commission_rule(scope=商品) / group_buy_activity
product_category 1─N commission_rule(scope=分类)
trade_order_item 1─1 product_review；1─N commission_record
trade_order 1─N group_buy_team_member (order_no)；N─1 promotion_activity / group_buy_team
coupon 1─N user_coupon 1─N invite_record(reward_ref)
commission_record 1─0..1 commission_record (冲销, reversal_record_id)
withdraw_order 1─N account_log；flash_sale_activity 1─N flash_sale_item
```

## 状态机汇总

| 状态机 | 状态与迁移 |
|---|---|
| notify_task.status | 10待发送→20已发送；→30失败→(重试)10 或 →40达上限；→50已跳过（渠道关闭） |
| distribution_user.status | 1待审核→2通过/3冻结（2↔3 可逆） |
| commission_record.status | 1待结算→2已结算/3已失效；2→4欠款冲销中（产生负额冲销记录回指原记录） |
| withdraw_order.status | 10待审核→20审核通过/50拒绝；20→30打款中；30→40成功/60失败已回退（条件更新幂等） |
| group_buy_team.status | 1拼团中→2已成团（人齐）/3已解散（超时，联动订单取消退款） |
| product_review.audit_status | 0待审→1通过/2驳回 |
| risk_record.appeal_status | 0无→1申诉中→2通过/3驳回 |
