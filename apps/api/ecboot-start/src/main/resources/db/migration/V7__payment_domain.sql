-- ============================================================
-- V7 支付域 (payment domain)
-- 表: pay_order(支付单), pay_callback_log(回调原文日志)
-- 一个订单可多次尝试支付(失败重发起新支付单),
-- 同一订单同时至多一个"待支付"支付单(服务层保证)
-- 回调幂等: 条件UPDATE抢占 + UNIQUE(渠道,渠道单号)防重放
-- 金额/币种回调时校验: amount=订单pay_amount, currency一致
-- ============================================================

CREATE TABLE `pay_order` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '支付单ID',
  `pay_no`           VARCHAR(32)     NOT NULL COMMENT '支付单号(全局唯一,传给支付渠道的商户单号)',
  `order_id`         BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no`         VARCHAR(32)     NOT NULL COMMENT '订单号',
  `user_id`          BIGINT UNSIGNED NOT NULL COMMENT '付款用户ID',
  `pay_channel`      TINYINT         NOT NULL DEFAULT 1 COMMENT '支付渠道:1微信支付 2支付宝(预留)',
  `amount`           DECIMAL(10,2)   NOT NULL COMMENT '应付金额(必须等于订单pay_amount,回调时校验)',
  `currency`         CHAR(3)         NOT NULL DEFAULT 'CNY' COMMENT '币种(ISO 4217,须与订单一致)',
  `status`           TINYINT         NOT NULL DEFAULT 10 COMMENT '支付状态:10待支付 20支付成功 30支付失败 90已关闭',
  `channel_trade_no` VARCHAR(64)     NULL DEFAULT NULL COMMENT '第三方交易号(微信/支付宝,成功回填;唯一防重放)',
  `success_time`     DATETIME(3)     NULL DEFAULT NULL COMMENT '支付成功时间(毫秒精度,渠道回调为准)',
  `expire_time`      DATETIME        NULL DEFAULT NULL COMMENT '支付截止时间(超时未付则关单)',
  `closed_time`      DATETIME        NULL DEFAULT NULL COMMENT '关闭时间',
  `fail_reason`      VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '失败原因(渠道返回)',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_pay_no` (`pay_no`),
  UNIQUE KEY `uk_channel_trade` (`pay_channel`, `channel_trade_no`),
  KEY `idx_order` (`order_id`, `status`),
  KEY `idx_user_time` (`user_id`, `created_at`),
  KEY `idx_status_expire` (`status`, `expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='支付单表(一次支付尝试一单)';

CREATE TABLE `pay_callback_log` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `pay_no`         VARCHAR(32)     NOT NULL COMMENT '关联支付单号',
  `pay_channel`    TINYINT         NOT NULL COMMENT '支付渠道:1微信 2支付宝',
  `notify_type`    TINYINT         NOT NULL DEFAULT 1 COMMENT '通知类型:1支付结果 2退款结果',
  `verify_status`  TINYINT         NOT NULL DEFAULT 0 COMMENT '验签结果:0失败 1通过',
  `process_status` TINYINT         NOT NULL DEFAULT 0 COMMENT '处理结果:0未处理/重复忽略 1已处理',
  `raw_body`       JSON            NOT NULL COMMENT '回调原文(全量留档,可重放审计)',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_pay_no` (`pay_no`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='支付回调日志(只追加,重复通知也留档)';
