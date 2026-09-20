-- ============================================================
-- V18 秒杀域 (flash sale)
-- 表: flash_sale_activity(活动), flash_sale_item(场次商品)
-- 活动库存与普通库存分账(research D7/契约R4):
--   可售 = stock_count - sold_count(推导值)
--   下单: UPDATE flash_sale_item SET sold_count=sold_count+N
--         WHERE id=? AND stock_count-sold_count>=N  (affected=0→已抢完)
--   取消: sold_count=sold_count-N(回补活动侧,不落inventory)
-- 秒杀价经 trade_order_item.price 快照承载;不超卖由条件更新保证(ADR-0001同构)
-- ============================================================

CREATE TABLE `flash_sale_activity` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `name`       VARCHAR(64)     NOT NULL COMMENT '活动名称',
  `start_time` DATETIME        NOT NULL COMMENT '开始时间(含)',
  `end_time`   DATETIME        NOT NULL COMMENT '结束时间(不含)',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status_time` (`status`, `start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='秒杀活动表(时段/状态)';

CREATE TABLE `flash_sale_item` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '场次商品ID',
  `activity_id`  BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `sku_id`       BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `flash_price`  DECIMAL(10,2)   NOT NULL COMMENT '秒杀价(快照进订单项price)',
  `stock_count`  INT UNSIGNED    NOT NULL COMMENT '活动限量(与inventory分账,契约R4)',
  `sold_count`   INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '已售(取消回补本列;可售=stock_count-sold_count)',
  `per_limit`    INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '每人限购',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_sku` (`activity_id`, `sku_id`),
  KEY `idx_sku` (`sku_id`),
  CONSTRAINT `chk_sold_le_stock` CHECK (`sold_count` <= `stock_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='秒杀场次商品表(活动分账库存,条件更新防超卖)';
