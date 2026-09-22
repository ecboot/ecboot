-- ============================================================
-- V42 撤销订单头 flash_sale_item_id 死列（015-marketing-c 评审修复 I6）
-- 更正: 000040 本可避免——trade_order_item.flash_sale_item_id 自 000023
--       就已存在且带 idx_flash_item（行级归属比订单头更自然）。本迁移
--       撤销订单头的同名列, 下单/取消回补/售后回补一律改用订单项既有列。
-- ============================================================
ALTER TABLE `trade_order`
  DROP KEY `idx_flash_item`,
  DROP COLUMN `flash_sale_item_id`;
