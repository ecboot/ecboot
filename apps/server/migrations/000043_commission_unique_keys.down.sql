-- 回滚 V43（N3 修复: 实测可逆形态——先删唯一键再删生成列, 顺序不可反）
ALTER TABLE `commission_record`
  DROP KEY `uk_item_bene_level_orig`,
  DROP KEY `uk_reversal_of`,
  DROP COLUMN `is_original`;
