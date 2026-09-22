-- ============================================================
-- V41 砍价单 order_no 改可空（015-marketing-c 评审修复 C1）
-- 缺陷: uk_order_no(order_no) + NOT NULL DEFAULT '' → 全库只能存在一张
--       砍价单（第二人发起必撞 1062）。"下单后回填"本就是 NULL 语义,
--       MySQL 唯一索引不约束 NULL → 未下单行写 NULL 即可共存任意多条。
-- ============================================================
UPDATE `bargain_record` SET `order_no` = NULL WHERE `order_no` = '';

ALTER TABLE `bargain_record`
  MODIFY COLUMN `order_no` VARCHAR(32) NULL DEFAULT NULL COMMENT '成交订单号(下单后回填; NULL=未下单, NULL不参与唯一约束)';
