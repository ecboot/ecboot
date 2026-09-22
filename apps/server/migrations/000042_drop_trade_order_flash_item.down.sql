ALTER TABLE `trade_order`
  ADD COLUMN `flash_sale_item_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '秒杀场次商品ID(NULL=普通单; 取消时据此回补活动库存)' AFTER `request_token`,
  ADD KEY `idx_flash_item` (`flash_sale_item_id`);
