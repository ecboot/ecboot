-- ============================================================
-- 000034 SPU 价格冗余列 (特性005, research D1)
-- 列表价格区间筛选与价格排序的索引友好实现:
--   SKU 创建/编辑/启停/删除时同事务重算启用 SKU 的 min/max 刷新本列
--   全部 SKU 禁用时为 NULL（卡片标不可下单）
-- ============================================================

ALTER TABLE `product_spu`
  ADD COLUMN `price_min` DECIMAL(10,2) NULL DEFAULT NULL COMMENT '启用SKU最低价(冗余,SKU变更时刷新;全部禁用为NULL)' AFTER `sale_count`,
  ADD COLUMN `price_max` DECIMAL(10,2) NULL DEFAULT NULL COMMENT '启用SKU最高价(冗余)' AFTER `price_min`,
  ADD INDEX `idx_price_min` (`price_min`);
