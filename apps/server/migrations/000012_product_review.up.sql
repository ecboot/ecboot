-- ============================================================
-- V12 评价域 (product review)
-- 表: product_review(商品评价)
-- 粒度: 订单项(一项一评, UNIQUE(order_item_id))
-- 快照: spu_name/sku_specs 冻结下单时商品信息(ADR-0002 原则,
--       商品软删后评价自洽展示, spec 边界场景)
-- 追评(extra_*)与商家回复(reply_*)各一次, 字段唯一承载
-- ============================================================

CREATE TABLE `product_review` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评价ID',
  `order_item_id`  BIGINT UNSIGNED NOT NULL COMMENT '订单项ID(一项一评)',
  `order_no`       VARCHAR(32)     NOT NULL COMMENT '订单号(冗余,免联查)',
  `user_id`        BIGINT UNSIGNED NOT NULL COMMENT '评价用户ID',
  `spu_id`         BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID(聚合统计用)',
  `sku_id`         BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `spu_name`       VARCHAR(128)    NOT NULL COMMENT 'SPU名称(下单时快照,商品软删后自洽展示)',
  `sku_specs`      JSON            NOT NULL COMMENT '规格组合(下单时快照)',
  `score`          TINYINT         NOT NULL COMMENT '评分:1-5星',
  `content`        VARCHAR(1024)   NOT NULL DEFAULT '' COMMENT '评价内容',
  `images`         JSON            NULL COMMENT '评价图片URL数组',
  `is_anonymous`   TINYINT         NOT NULL DEFAULT 0 COMMENT '匿名:0否 1是',
  `audit_status`   TINYINT         NOT NULL DEFAULT 0 COMMENT '审核状态:0待审核 1通过 2驳回',
  `extra_content`  VARCHAR(1024)   NOT NULL DEFAULT '' COMMENT '追评内容(仅一次,90天内)',
  `extra_time`     DATETIME        NULL DEFAULT NULL COMMENT '追评时间',
  `reply_content`  VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '商家回复(仅一次)',
  `reply_time`     DATETIME        NULL DEFAULT NULL COMMENT '回复时间',
  `deleted`        TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是(用户删除自己的评价)',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_item` (`order_item_id`),
  KEY `idx_spu_audit_time` (`spu_id`, `audit_status`, `created_at`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品评价表(订单项粒度,审核后展示)';
