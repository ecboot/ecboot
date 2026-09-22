-- 回滚 V41（018 批次 13 收口终验 C2: 逆向顺序修正 + 多 NULL 处理）
-- 不可直接 UPDATE 置 '' 再改 NOT NULL: 多行 NULL 会撞 uk_order_no（Error 1062 Duplicate entry ''）.
-- 实际语义（N7 二次收口——原注释声称"多行时报错中止"与 SQL 行为不符, 已按实际行为改写）:
--   * NULL 行 ≤ 1 → 直接回填空串, 语义无损;
--   * NULL 行 ≥ 2 → **有损回滚**: 逐行回填 'LEGACY-<id>' 唯一占位（保 uk_order_no 不撞,
--     但这些行回滚后形如"已下单"——需人工核对; 属不可完全逆的语义降级, 已在 §六 登记）。
SET @n := (SELECT COUNT(*) FROM `bargain_record` WHERE `order_no` IS NULL);
UPDATE `bargain_record`
   SET `order_no` = CONCAT('LEGACY-', `id`)
 WHERE `order_no` IS NULL AND @n > 1;

UPDATE `bargain_record` SET `order_no` = '' WHERE `order_no` IS NULL;

ALTER TABLE `bargain_record`
  MODIFY COLUMN `order_no` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '成交订单号(下单后回填,防重复成交)';
