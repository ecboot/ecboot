-- ============================================================
-- V17 拼团域 (group buy)
-- 表: group_buy_activity(活动), group_buy_team(团),
--     group_buy_team_member(团员) + trade_order 增列
-- 约束: UNIQUE(team_id,user_id)=一团一单; UNIQUE(order_no)=一单一团
-- 成团价经 trade_order_item.price 快照承载(research D10)
-- 超时未成团→系统取消订单+全额退款(应用层编排,表仅承载状态)
-- ============================================================

CREATE TABLE `group_buy_activity` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `spu_id`         BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID(成团价对全SKU生效)',
  `name`           VARCHAR(64)     NOT NULL COMMENT '活动名称',
  `group_price`    DECIMAL(10,2)   NOT NULL COMMENT '成团价(下单快照进订单项price)',
  `group_size`     INT UNSIGNED    NOT NULL DEFAULT 2 COMMENT '成团人数',
  `valid_start_at` DATETIME        NOT NULL COMMENT '活动开始时间',
  `valid_end_at`   DATETIME        NOT NULL COMMENT '活动结束时间',
  `per_limit`      INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '每人限购(整个活动期)',
  `status`         TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`        TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_spu_status` (`spu_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='拼团活动表(成团人数/时限/限购)';

CREATE TABLE `group_buy_team` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '团ID',
  `activity_id`    BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `leader_user_id` BIGINT UNSIGNED NOT NULL COMMENT '团长用户ID',
  `status`         TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1拼团中 2已成团 3已解散(超时/取消)',
  `member_count`   INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '当前成员数(应用层与member表同事务维护)',
  `expire_time`    DATETIME        NOT NULL COMMENT '成团截止时间(超时扫描依据)',
  `success_time`   DATETIME        NULL DEFAULT NULL COMMENT '成团时间',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开团时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_activity_status` (`activity_id`, `status`),
  KEY `idx_status_expire` (`status`, `expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='拼团团表(人齐成团/超时解散)';

CREATE TABLE `group_buy_team_member` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '团员ID',
  `team_id`    BIGINT UNSIGNED NOT NULL COMMENT '团ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '成员用户ID',
  `order_no`   VARCHAR(32)     NOT NULL COMMENT '成员订单号',
  `join_time`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '参团时间',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_team_user` (`team_id`, `user_id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='拼团团员表(一团一单/一单一团,双重唯一)';

-- 订单关联列(可空,普通订单零影响)
ALTER TABLE `trade_order`
  ADD COLUMN `group_buy_team_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '拼团团ID(NULL=普通订单)' AFTER `user_coupon_id`,
  ADD KEY `idx_group_team` (`group_buy_team_id`);
