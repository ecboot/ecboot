-- ============================================================
-- V24 拼团定价粒度下沉 SKU 级 (产品决策 2026-09-18)
-- 评审遗留项1落地: 对齐 flash_sale_item 模式——多规格商品
-- 按规格差异化定价, 消除 SPU 级单值的表达局限
-- 库存策略不变: 拼团走普通库存(下单锁定制), 无活动分账
-- ============================================================

CREATE TABLE `group_buy_item` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '拼团场次商品ID',
  `activity_id` BIGINT UNSIGNED NOT NULL COMMENT '拼团活动ID',
  `sku_id`      BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `group_price` DECIMAL(10,2)   NOT NULL COMMENT '该SKU成团价(下单快照进订单项price)',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_sku` (`activity_id`, `sku_id`),
  KEY `idx_sku` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='拼团场次商品表(SKU级成团价;限购per_limit保持活动级;库存走普通inventory)';

-- 移除活动级单值价格, 避免双真源(库内无数据, 无迁移步骤)
ALTER TABLE `group_buy_activity` DROP COLUMN `group_price`;
