-- 000038_after_sale_operator 回退
ALTER TABLE `after_sale_order`
  DROP COLUMN `fail_reason`,
  DROP COLUMN `operator_id`;
