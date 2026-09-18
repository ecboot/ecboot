-- ============================================================
-- V4 购物车域 (cart domain)
-- 表: cart_item(购物车项)
-- 价格/库存不快照: 购物车实时联查SKU(价格变动以结算页为准)
-- 重复加购 = 更新数量 (UNIQUE(user_id, sku_id))
-- 移除 = 物理删除 (无deleted字段)
-- ============================================================

CREATE TABLE `cart_item` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '购物车项ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `sku_id`     BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `quantity`   INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '数量(应用层限制上限,如99)',
  `checked`    TINYINT         NOT NULL DEFAULT 1 COMMENT '结算勾选:0未选 1选中',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_sku` (`user_id`, `sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='购物车表(有效SKU由结算时校验,下架/禁用商品置灰)';
