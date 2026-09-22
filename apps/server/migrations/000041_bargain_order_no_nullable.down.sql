UPDATE `bargain_record` SET `order_no` = '' WHERE `order_no` IS NULL;

ALTER TABLE `bargain_record`
  MODIFY COLUMN `order_no` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '成交订单号(下单后回填,防重复成交)';
