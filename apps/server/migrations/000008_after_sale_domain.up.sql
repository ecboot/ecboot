-- ============================================================
-- V8 售后域 (after-sale domain)
-- 表: after_sale_order(售后单)
-- 类型: 1仅退款 2退货退款 (3换货预留值,不启用)
-- 粒度: 按订单项发起(支持部分退款),一单可多次售后(不同项)
-- 金额: refund_amount ≤ 订单项pay_amount,应用层校验;币种随订单
-- 渠道退款幂等: out_refund_no = after_sale_no(渠道侧幂等键)
-- ============================================================

CREATE TABLE `after_sale_order` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '售后单ID',
  `after_sale_no`       VARCHAR(32)     NOT NULL COMMENT '售后单号(全局唯一,兼作渠道退款out_refund_no幂等键)',
  `order_id`            BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no`            VARCHAR(32)     NOT NULL COMMENT '订单号',
  `order_item_id`       BIGINT UNSIGNED NOT NULL COMMENT '订单项ID(售后粒度=订单项)',
  `user_id`             BIGINT UNSIGNED NOT NULL COMMENT '申请人用户ID',
  `type`                TINYINT         NOT NULL COMMENT '售后类型:1仅退款 2退货退款 3换货(预留未启用)',
  `status`              TINYINT         NOT NULL DEFAULT 10 COMMENT '状态:10待审核 20待买家寄回 30待退款 40退款中 50已完成 90已拒绝 91已撤销',
  `currency`            CHAR(3)         NOT NULL DEFAULT 'CNY' COMMENT '币种(ISO 4217,随订单)',
  `reason`              VARCHAR(64)     NOT NULL COMMENT '售后原因(不想要了/质量问题/发错货等)',
  `description`         VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '问题描述',
  `voucher_images`      JSON            NULL COMMENT '凭证图片URL数组',
  `refund_amount`       DECIMAL(10,2)   NOT NULL COMMENT '退款金额(≤订单项pay_amount)',
  `return_logistics_no` VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '买家寄回物流单号(type=2)',
  `refund_no`           VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '渠道退款单号(微信退款ID,回填)',
  `reject_reason`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '拒绝原因',
  `audit_time`          DATETIME        NULL DEFAULT NULL COMMENT '审核时间',
  `refund_time`         DATETIME        NULL DEFAULT NULL COMMENT '退款完成时间',
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_after_sale_no` (`after_sale_no`),
  KEY `idx_order` (`order_id`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='售后单表(退款状态机载体,完成时回补库存并更新订单refund_status)';
