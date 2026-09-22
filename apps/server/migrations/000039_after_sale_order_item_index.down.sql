-- 000039_after_sale_order_item_index 回退
ALTER TABLE `after_sale_order`
  DROP KEY `idx_order_item`;
