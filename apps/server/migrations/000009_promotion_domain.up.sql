-- ============================================================
-- V9 促销域 (promotion domain) — V1 仅优惠券
-- 表: coupon(优惠券模板), user_coupon(用户持有的券)
-- 券类型: 1满减券(满X减Y) 2无门槛券 3折扣券(预留未启用)
-- 币种: 优惠券仅人民币CNY(多币种抵扣涉及汇率折算,明确不支持)
-- 核销: 下单成功同事务置"已使用",订单取消置"已退回"
-- 防超发: UPDATE coupon SET received_count=received_count+1
--         WHERE id=? AND (total_count=0 OR received_count<total_count)
-- ============================================================

CREATE TABLE `coupon` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '优惠券模板ID',
  `name`              VARCHAR(64)     NOT NULL COMMENT '券名称(如:满100减20)',
  `type`              TINYINT         NOT NULL DEFAULT 1 COMMENT '券类型:1满减券 2无门槛券 3折扣券(预留)',
  `threshold_amount`  DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '使用门槛(满X元可用,0=无门槛)',
  `discount_amount`   DECIMAL(10,2)   NOT NULL COMMENT '抵扣金额(满减/无门槛券)',
  `discount_rate`     DECIMAL(3,2)    NULL DEFAULT NULL COMMENT '折扣率0.01-0.99(折扣券预留,如0.85=85折)',
  `total_count`       INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '发放总量,0=不限量',
  `received_count`    INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '已领取数量(领券事务内条件更新,防超发)',
  `per_limit`         INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '每人限领数量',
  `valid_type`        TINYINT         NOT NULL DEFAULT 1 COMMENT '有效期方式:1固定区间 2领取后N天',
  `valid_start_at`    DATETIME        NULL DEFAULT NULL COMMENT '固定区间开始(valid_type=1)',
  `valid_end_at`      DATETIME        NULL DEFAULT NULL COMMENT '固定区间结束(valid_type=1)',
  `valid_days`        INT UNSIGNED    NULL DEFAULT NULL COMMENT '领取后N天有效(valid_type=2)',
  `status`            TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用(停用不影响已领取的券)',
  `deleted`           TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='优惠券模板表(V1全场通用;品类/商品限定范围预留)';

CREATE TABLE `user_coupon` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户券ID(订单user_coupon_id引用此表)',
  `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '持券用户ID',
  `coupon_id`   BIGINT UNSIGNED NOT NULL COMMENT '优惠券模板ID',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1未使用 2已使用 3已过期 4已退回(订单取消,可再用)',
  `expire_time` DATETIME        NOT NULL COMMENT '过期时间(领取时按模板规则计算落库)',
  `order_no`    VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '核销订单号',
  `used_time`   DATETIME        NULL DEFAULT NULL COMMENT '核销时间',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_status_expire` (`user_id`, `status`, `expire_time`),
  KEY `idx_coupon` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户持有的优惠券(过期由过期任务或使用时惰性判定置3)';
