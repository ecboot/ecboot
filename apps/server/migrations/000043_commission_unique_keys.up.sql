-- ============================================================
-- V43 佣金记录补唯一键（017 批次 11 评审修复 C3/C4——零迁移承诺按协议破除记账）
-- 缺陷: 计提/冲销均为 check-then-insert 且无唯一索引 → 并发双计提/双冲销
--       （评审探针: 双计提 rows=2; 双扣 100→80 应 90）。
-- 设计: MySQL 唯一索引含 NULL 的行不参与唯一性检查——
--   is_original 生成列: 原始行(未冲销)=1 参与唯一 → uk_item_bene_level_orig 保证
--   (order_item_id,beneficiary,level) 只一条原始行（计提幂等）;
--   冲销行 is_original=NULL 不参与该键（由 uk_reversal_of 一对一保证冲销幂等）。
-- ============================================================
ALTER TABLE `commission_record`
  ADD COLUMN `is_original` TINYINT GENERATED ALWAYS AS (IF(`reversal_of_id` IS NULL, 1, NULL)) STORED,
  ADD UNIQUE KEY `uk_item_bene_level_orig` (`order_item_id`, `beneficiary_user_id`, `level`, `is_original`),
  ADD UNIQUE KEY `uk_reversal_of` (`reversal_of_id`);
