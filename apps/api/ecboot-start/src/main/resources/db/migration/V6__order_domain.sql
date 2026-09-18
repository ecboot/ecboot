-- ============================================================
-- V6 订单域 (order/trade domain)
-- 表: trade_order(订单), trade_order_item(订单项,含快照),
--     trade_order_log(订单状态流水)
-- 状态机: 10待付款→20待发货→30待收货→40已完成
--         10→90已取消(用户/30分钟超时/管理员)
-- 幂等: UNIQUE(user_id, request_token) 兜底下单防重
-- 币种: currency CHAR(3) ISO 4217,默认CNY;明细继承头表
-- 运费: 按运费模板计算,下单时快照(freight_amount+freight_template_id)
-- 交易单据永不软删(取消是状态迁移,不是删除)
-- ============================================================

CREATE TABLE `trade_order` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no`            VARCHAR(32)     NOT NULL COMMENT '订单号(业务号:日期+雪花/随机,全局唯一,分片友好)',
  `user_id`             BIGINT UNSIGNED NOT NULL COMMENT '买家用户ID',
  `order_channel`       TINYINT         NOT NULL DEFAULT 1 COMMENT '下单渠道:1微信小程序 2H5',
  `status`              TINYINT         NOT NULL DEFAULT 10 COMMENT '订单状态:10待付款 20待发货 30待收货 40已完成 90已取消',
  `refund_status`       TINYINT         NOT NULL DEFAULT 0 COMMENT '退款状态:0无售后 1部分退款 2全额退款(不影响主状态机)',
  `currency`            CHAR(3)         NOT NULL DEFAULT 'CNY' COMMENT '币种(ISO 4217,订单项继承此字段)',
  `total_amount`        DECIMAL(10,2)   NOT NULL COMMENT '商品总额(原价*数量合计)',
  `promotion_amount`    DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '优惠总额(优惠券等)',
  `freight_amount`      DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '运费(按运费模板计算,下单时快照)',
  `freight_template_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '运费模板ID(下单时快照,NULL=包邮)',
  `pay_amount`          DECIMAL(10,2)   NOT NULL COMMENT '实付金额=total_amount-promotion_amount+freight_amount',
  `user_coupon_id`      BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '使用的用户优惠券ID(与订单同事务核销)',
  `receiver_name`       VARCHAR(64)     NOT NULL COMMENT '收货人姓名(下单时快照)',
  `receiver_phone`      VARCHAR(20)     NOT NULL COMMENT '收货人手机号(快照)',
  `receiver_province`   VARCHAR(32)     NOT NULL COMMENT '省(快照,运费区域匹配依据)',
  `receiver_city`       VARCHAR(32)     NOT NULL COMMENT '市(快照)',
  `receiver_district`   VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '区/县(快照)',
  `receiver_detail`     VARCHAR(255)    NOT NULL COMMENT '详细地址(快照)',
  `user_remark`         VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '买家留言',
  `request_token`       VARCHAR(64)     NULL DEFAULT NULL COMMENT '下单幂等token(确认页发放,Redis抢占+唯一索引兜底;NULL不参与唯一)',
  `pay_time`            DATETIME        NULL DEFAULT NULL COMMENT '支付完成时间',
  `deliver_company`     VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '物流公司(发货预留)',
  `deliver_no`          VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '物流单号(发货预留)',
  `deliver_time`        DATETIME        NULL DEFAULT NULL COMMENT '发货时间(预留)',
  `finish_time`         DATETIME        NULL DEFAULT NULL COMMENT '订单完成时间(确认收货)',
  `cancel_type`         TINYINT         NULL DEFAULT NULL COMMENT '取消方:1用户 2系统超时 3管理员',
  `cancel_reason`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '取消原因',
  `cancel_time`         DATETIME        NULL DEFAULT NULL COMMENT '取消时间',
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  UNIQUE KEY `uk_user_token` (`user_id`, `request_token`),
  KEY `idx_user_status_time` (`user_id`, `status`, `created_at`),
  KEY `idx_status_time` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单主表';

CREATE TABLE `trade_order_item` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单项ID',
  `order_id`          BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no`          VARCHAR(32)     NOT NULL COMMENT '订单号(冗余,免联查)',
  `spu_id`            BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID(溯源用)',
  `sku_id`            BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID(售后/库存回补定位)',
  `sku_no`            VARCHAR(32)     NOT NULL COMMENT 'SKU编码(快照)',
  `spu_name`          VARCHAR(128)    NOT NULL COMMENT 'SPU名称(下单时快照)',
  `sku_name`          VARCHAR(128)    NOT NULL COMMENT 'SKU名称=名称+规格串(快照)',
  `sku_image`         VARCHAR(512)    NOT NULL DEFAULT '' COMMENT 'SKU主图URL(下单时快照)',
  `sku_specs`         JSON            NOT NULL COMMENT '规格组合(快照):{"颜色":"黑","尺码":"M"}',
  `quantity`          INT UNSIGNED    NOT NULL COMMENT '购买数量',
  `original_price`    DECIMAL(10,2)   NOT NULL COMMENT '下单时页面价(快照)',
  `price`             DECIMAL(10,2)   NOT NULL COMMENT '成交单价(优惠分摊前)',
  `promotion_amount`  DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '分摊到本行的优惠(按行金额比例,尾差记末行)',
  `pay_amount`        DECIMAL(10,2)   NOT NULL COMMENT '本行实付=price*quantity-promotion_amount(币种继承订单头表)',
  `created_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间(与created_at相等,订单项不可变)',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_sku` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单项表(商品快照,售后金额依据)';

CREATE TABLE `trade_order_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `order_id`      BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no`      VARCHAR(32)     NOT NULL COMMENT '订单号',
  `from_status`   TINYINT         NULL DEFAULT NULL COMMENT '迁移前状态(建单时为NULL)',
  `to_status`     TINYINT         NOT NULL COMMENT '迁移后状态',
  `operator_type` TINYINT         NOT NULL DEFAULT 1 COMMENT '操作者类型:1系统 2用户 3管理员',
  `operator_id`   VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '操作者标识(system/user:{id}/admin:{id})',
  `remark`        VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '备注(如超时取消/确认收货)',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单状态迁移流水(只追加,排查与审计)';
