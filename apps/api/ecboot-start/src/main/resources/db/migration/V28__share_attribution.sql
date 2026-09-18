-- ============================================================
-- V28 分享归因域 (社交电商第二轮评审 P0+P1, 产品决策 2026-09-18)
-- 归因模型(佣金计提判定依据,三级优先):
--   1 分享归因: 订单来自 share_record 的分享触点(默认窗口7天,
--     应用配置化), 归因人=分享人 → 谁分享谁受益
--   2 关系链兜底: 无归因时按 user_relation 两级计佣
--   3 自然流量: 无归因无关系链, 不计佣
-- 规则: 自购不算自己的归因(防自刷; 自购返佣为独立规则,待决策)
-- 附: 首单激励时机(P1-3) + 卖家备注(P1-4)
-- ============================================================

-- 个人推广码(海报/口令的个人标识;生成策略应用层——建议注册即生成)
ALTER TABLE `user`
  ADD COLUMN `share_code` VARCHAR(16) NULL DEFAULT NULL COMMENT '个人推广码(唯一;NULL=未生成,唯一索引不参与)' AFTER `wx_unionid`,
  ADD UNIQUE KEY `uk_share_code` (`share_code`);

-- 分享行为记录(只追加;归因判定入口+分享行为分析)
CREATE TABLE `share_record` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分享记录ID',
  `sharer_user_id` BIGINT UNSIGNED NOT NULL COMMENT '分享人用户ID',
  `spu_id`        BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '分享对象SPU(NULL=整店/页面分享)',
  `share_channel` TINYINT         NOT NULL COMMENT '渠道:1小程序卡片 2海报图片 3复制链接/口令 4朋友圈/社群',
  `scene_value`   VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '小程序场景值(scene参数,落地追踪)',
  `share_code`    VARCHAR(16)     NOT NULL DEFAULT '' COMMENT '使用的推广码',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分享时间(归因窗口起点)',
  PRIMARY KEY (`id`),
  KEY `idx_sharer_time` (`sharer_user_id`, `created_at`),
  KEY `idx_spu_time` (`spu_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='分享行为记录表(只追加,订单归因的判定来源)';

-- 订单归因(佣金计提:归因优先于关系链)
ALTER TABLE `trade_order`
  ADD COLUMN `attributed_user_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '订单归因人(分享带来本单的用户;佣金按 分享归因>关系链兜底>自然流量 三级判定)' AFTER `user_coupon_id`,
  ADD COLUMN `attribution_type`   TINYINT         NOT NULL DEFAULT 3 COMMENT '归因类型:1分享归因 2关系链兜底 3自然流量' AFTER `attributed_user_id`,
  ADD COLUMN `seller_remark`      VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '卖家备注(客服/仓库内部使用,买家不可见)' AFTER `user_remark`,
  ADD KEY `idx_attributed` (`attributed_user_id`);

-- 首单激励时机(邀请奖励从"注册发放"扩展为可选"首单发放",防刷主流形态)
ALTER TABLE `invite_record`
  ADD COLUMN `reward_trigger` TINYINT NOT NULL DEFAULT 1 COMMENT '奖励触发时机:1注册即发 2首单后发(防刷,需订单回调触发)' AFTER `reward_type`;
