-- 000037_user_level_bigint 回退
-- 警告: 若库内已有 > 127 的等级值, 本回退会报 1264 或截断数据 —— 回退前须先人工核处这些行
-- （`SELECT COUNT(*) FROM user WHERE level > 127`）。
ALTER TABLE `user`
  MODIFY COLUMN `level` TINYINT NULL DEFAULT NULL COMMENT '会员等级(user_level_rule.id,未启体系为NULL)';
