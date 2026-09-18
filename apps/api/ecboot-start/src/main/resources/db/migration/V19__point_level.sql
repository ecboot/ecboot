-- ============================================================
-- V19 积分与会员等级域 (points & level)
-- 表: point_account, point_log, user_level_rule
--     + user 增列 + trade_order/trade_order_item 积分列
-- 双账本分离(clarify Q3/research D5):
--   积分=可消耗(账户+流水,余额可负=欠款抵扣)
--   成长值=只增不减(user.growth_value 累计列,退款不扣)
-- 优惠三构成之一: point_amount(研究D8,恒等式见V20与quickstart场景三)
-- ============================================================

CREATE TABLE `point_account` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '积分账户ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `balance`    INT             NOT NULL DEFAULT 0 COMMENT '积分余额(有符号,回退欠款为负,禁止UNSIGNED)',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='积分账户表(可负余额)';

CREATE TABLE `point_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `biz_type`      TINYINT         NOT NULL COMMENT '业务类型:1签到 2消费获得 3下单消耗 4退款回退(负) 5分享获得',
  `points`        INT             NOT NULL COMMENT '积分变动(有符号:获得正,消耗/回退负)',
  `balance_after` INT             NOT NULL COMMENT '变动后余额快照(对账用)',
  `order_no`      VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '关联订单号',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`, `id`),
  KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='积分流水表(只追加,获取/消耗/回退全留痕)';

CREATE TABLE `user_level_rule` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '等级规则ID',
  `name`             VARCHAR(32)     NOT NULL COMMENT '等级名称(如 普通会员/白银/黄金)',
  `growth_threshold` INT UNSIGNED    NOT NULL COMMENT '成长值门槛(唯一,等级判定=满足的最高门槛)',
  `benefits`         JSON            NULL COMMENT '等级权益配置(折扣/优先客服等,键值对)',
  `status`           TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`          TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_threshold` (`growth_threshold`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员等级规则表(成长值门槛唯一递增)';

-- 成长值与等级(等级可空:未启等级体系时零影响,FR-019)
ALTER TABLE `user`
  ADD COLUMN `growth_value` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '成长值(只增不减,累计列)' AFTER `gender`,
  ADD COLUMN `level` TINYINT NULL DEFAULT NULL COMMENT '会员等级(user_level_rule.id,未启体系为NULL)' AFTER `growth_value`;

-- 订单层积分抵扣(优惠三构成之一)
ALTER TABLE `trade_order`
  ADD COLUMN `point_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '积分抵扣金额(promotion_amount构成之一)' AFTER `freight_template_id`,
  ADD COLUMN `point_used` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '消耗积分点数' AFTER `point_amount`;

ALTER TABLE `trade_order_item`
  ADD COLUMN `point_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '积分抵扣行分摊(合计=订单头point_amount,尾差记末行)' AFTER `promotion_amount`;
