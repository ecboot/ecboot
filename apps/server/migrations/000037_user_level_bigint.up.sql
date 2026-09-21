-- 000037_user_level_bigint: user.level 列类型对齐其自身契约（012 评审修复轮发现并记账）
-- 契约: 该列存 **user_level_rule.id**（000019 建列注释 与 生成物 entity 注释 均如此声明）,
--       而被引用的 user_level_rule.id 是 BIGINT UNSIGNED。
-- 缺陷: 原列类型 TINYINT(±127) 与契约矛盾 —— 规则表 id > 127 时, LevelRecalc 写 user.level 报
--       `1264 Out of range value`, 该错误从 GrowthAdd 冒泡出去, **会员成长/升级链路直接失败**
--       （测试库反复重建规则把 AUTO_INCREMENT 推过 127 后必然复现）。
-- 语义不变: 仍存规则 ID, 仍可为空（NULL = 未启等级体系）; 仅放宽值域以匹配被引用列。
-- 注: 该列目前是"只写不读"的反规范化缓存（ProfileDetail 用 matchLevel 实时重算等级名）,
--     故本次放宽对读取路径零影响。
ALTER TABLE `user`
  MODIFY COLUMN `level` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '会员等级(user_level_rule.id,未启体系为NULL)';
