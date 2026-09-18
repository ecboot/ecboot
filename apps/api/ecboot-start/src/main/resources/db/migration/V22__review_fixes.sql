-- ============================================================
-- V22 评审修复 (code review fixes, b141db4 评审)
-- 1) 满减范围"每活动至多一条全场行"数据库级强制(评审 I2):
--    原唯一键 (activity_id, scope_type, target_id) 因 NULL≠NULL
--    无法约束 target_id=NULL 的全场行重复;
--    生成列 target_id_norm=IFNULL(target_id,0) + 新唯一键替换原键
-- 2) 关联列索引补齐(评审 I6 + Minor):
--    trade_order.promotion_activity_id / product_review.order_no+sku_id /
--    user_message.biz_no / group_buy_team.leader_user_id /
--    commission_record 冲销双向列 / user_footprint.spu_id
--    (FR-001: 关联列一律建索引;已应用迁移不可回改,故以 V22 增量补齐)
-- ============================================================

ALTER TABLE `promotion_activity_scope`
  ADD COLUMN `target_id_norm` BIGINT AS (IFNULL(`target_id`, 0)) STORED COMMENT '唯一键载体:全场(NULL)归一为0,堵住NULL≠NULL多行漏洞' AFTER `target_id`,
  ADD UNIQUE KEY `uk_activity_scope_norm` (`activity_id`, `scope_type`, `target_id_norm`),
  DROP INDEX `uk_activity_scope`;

ALTER TABLE `trade_order`
  ADD KEY `idx_promotion_activity` (`promotion_activity_id`);

ALTER TABLE `product_review`
  ADD KEY `idx_order_no` (`order_no`),
  ADD KEY `idx_sku` (`sku_id`);

ALTER TABLE `user_message`
  ADD KEY `idx_biz_no` (`biz_no`);

ALTER TABLE `group_buy_team`
  ADD KEY `idx_leader` (`leader_user_id`);

ALTER TABLE `commission_record`
  ADD KEY `idx_reversal_record` (`reversal_record_id`),
  ADD KEY `idx_reversal_of` (`reversal_of_id`);

ALTER TABLE `user_footprint`
  ADD KEY `idx_spu` (`spu_id`);
