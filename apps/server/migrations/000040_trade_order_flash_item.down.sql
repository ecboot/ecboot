-- 000040_trade_order_flash_item 回退
ALTER TABLE `trade_order`
  DROP KEY `idx_flash_item`,
  DROP COLUMN `flash_sale_item_id`;
