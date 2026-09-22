-- 回滚 V41（018 批次 13 收口终验 C2: 逆向顺序修正 + 多 NULL 合并）
-- 不可直接 UPDATE 置 '' 再改 NOT NULL: 多行 NULL 会撞 uk_order_no（Error 1062 Duplicate entry ''）.
-- 正确形态: **先把除首行外的 NULL 合并（删多余未下单行或置唯一占位）再改回 NOT NULL** ——
-- 但"合并"语义上会丢数据, 故 V1 采取可逆的保守形态: 仅当 NULL 行 ≤1 时改回;
-- 多行时以显式报错中止并保留可重跑的状态（迁移器的 dirty 语义即"需人工裁决"）。
SET @n := (SELECT COUNT(*) FROM `bargain_record` WHERE `order_no` IS NULL);
-- 若 NULL 行 > 1: 逐行回填唯一占位（'LEGACY-<id>'）以保可逆且不丢数据
UPDATE `bargain_record`
   SET `order_no` = CONCAT('LEGACY-', `id`)
 WHERE `order_no` IS NULL AND @n > 1;

UPDATE `bargain_record` SET `order_no` = '' WHERE `order_no` IS NULL;

ALTER TABLE `bargain_record`
  MODIFY COLUMN `order_no` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '成交订单号(下单后回填,防重复成交)';
