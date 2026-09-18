-- ============================================================
-- V26 运营与配置域 (产品决策落地 2026-09-18, 评审遗留 2/3/4/5)
-- ① 地区限售: SPU 级黑名单(JSON, 下单单点校验; 列表过滤交ES)
-- ② 积分滚动过期: 账户级"最后获得日+12个月", 批次制为升级路径
--    (触发条件: 积分商城/兑换比例变化/财务审计要求)
-- ③ 通知模板: 渠道现实建模——平台侧模板ID+参数契约+站内信全文自控
-- ④ 运营位: 轻量两表(banner+floor), 不做装修系统(素材直存OSS URL)
-- 账号归并(遗留1): 无DDL, 登录流程规则——微信登录取手机号→phone_hash
--    命中旧账号则绑定openid直接登录(冲突消解在入口), 不做迁移式合并
-- ============================================================

-- ① 地区限售(SPU级黑名单)
ALTER TABLE `product_spu`
  ADD COLUMN `sale_restrict_codes` JSON NULL COMMENT '禁售省级区划代码列表(黑名单,NULL=全国可售;下单时地址省代码∈列表即拒绝;列表页过滤由ES承载)' AFTER `attributes`;

-- ② 积分滚动过期(账户级口径)
ALTER TABLE `point_account`
  ADD COLUMN `last_earned_at` DATETIME NULL COMMENT '最后获得时间(滚动有效期口径:自最后获得日起12个月内有效;过期任务清零并记流水)' AFTER `balance`;

ALTER TABLE `point_log`
  MODIFY COLUMN `biz_type` TINYINT NOT NULL COMMENT '业务类型:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励 9过期扣减';

-- ③ 通知模板(运营可配;微信/短信模板在平台侧申请审核,本表管理映射与契约)
CREATE TABLE `notify_template` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `code`                 VARCHAR(64)     NOT NULL COMMENT '模板编码(notify_task.template_code引用)',
  `channel`              TINYINT         NOT NULL COMMENT '渠道:1小程序订阅消息 2短信 3站内信',
  `external_template_id` VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '平台侧模板ID(微信订阅消息/短信平台模板,平台审核产物;站内信为空)',
  `title`                VARCHAR(128)    NOT NULL DEFAULT '' COMMENT '标题(站内信自控;渠道侧为文案备份)',
  `content_template`     VARCHAR(2048)   NOT NULL COMMENT '内容模板({{变量}}占位;站内信全文自控,渠道侧为文案备份)',
  `params_schema`        JSON            NULL COMMENT '变量契约:[{"name":"orderNo","type":"string","example":"O1"}]开发与运营的参数约定',
  `status`               TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用(渠道级开关,短信成本闸门)',
  `deleted`              TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`           DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`           DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_channel` (`code`, `channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知模板表(一code一渠道;改文案不发版)';

-- ④ 运营位(轻量两表;小程序审核周期1-3天 vs 运营节奏矛盾的解法)
CREATE TABLE `operation_banner` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'bannerID',
  `position`   TINYINT         NOT NULL DEFAULT 1 COMMENT '位置:1首页轮播 2首页弹窗(枚举可扩展)',
  `image_url`  VARCHAR(512)    NOT NULL COMMENT '图片URL(OSS/CDN)',
  `link_url`   VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '跳转链接(小程序页面路径或H5;空=纯展示)',
  `sort`       INT             NOT NULL DEFAULT 0 COMMENT '排序,越小越靠前',
  `start_time` DATETIME        NULL DEFAULT NULL COMMENT '投放开始(NULL=立即生效)',
  `end_time`   DATETIME        NULL DEFAULT NULL COMMENT '投放结束(NULL=长期有效)',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_position_status` (`position`, `status`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='运营banner表(轮播/弹窗,投放时段可配)';

CREATE TABLE `operation_floor` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '楼层ID',
  `floor_type` TINYINT         NOT NULL COMMENT '类型:1金刚区 2商品楼层 3专题',
  `title`      VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '楼层标题',
  `config`     JSON            NULL COMMENT '楼层配置(金刚区=图标入口数组/商品楼层=商品ID列表与参数;类型化schema由应用层约定)',
  `sort`       INT             NOT NULL DEFAULT 0 COMMENT '排序,越小越靠前',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status_sort` (`status`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='运营楼层表(首页内容配置,改版不发版)';
