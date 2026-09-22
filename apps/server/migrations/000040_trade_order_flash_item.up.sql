-- 000040_trade_order_flash_item: 订单记录所属秒杀场次商品（015-marketing-c 批次 09）
-- 背景: 秒杀订单取消须**回补活动库存**（`flash_sale_item.sold_count -= 数量`）——这是批次 07 记下的
--   秒杀欠账清偿条件之一。而 `trade_order` 未记录该单来自哪个场次商品; 若改按 `sku_id` 反查"进行中的场次",
--   在"活动已结束才取消"时会查不到 → 漏回补 → 活动库存永久泄漏。
-- 语义: NULL = 普通单; 非空 = 该单使用该场次商品的秒杀价（下单时快照该 id, 取消时据此回补）。
ALTER TABLE `trade_order`
  ADD COLUMN `flash_sale_item_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '秒杀场次商品ID(NULL=普通单; 取消时据此回补活动库存)' AFTER `request_token`,
  ADD KEY `idx_flash_item` (`flash_sale_item_id`);
