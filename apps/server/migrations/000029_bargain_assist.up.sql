-- ============================================================
-- V29 砍价与助力玩法 + 分销规则落地 (产品决策 2026-09-18)
-- ① 自购返佣(规则级): 买家为通过审核的推广员时,一级佣金
--    受益人=本人("自买自省"),其关系链上级二级照常——与正常
--    订单同构,无 DDL;自购仍不算自己的"归因"(防自刷,归因与
--    佣金是两层规则)
-- ② 推广员等级: distribution_user.level 预留(V1 单一等级;
--    推广员>500人时启用多级比例,届时 commission_rule 加等级维度)
-- ③ 砍价(bargain): SKU级价格区间,好友各砍一刀(随机递减),
--    到底价可下单(订单价快照),超时失败——三段式对齐拼团/秒杀
-- ④ 助力(assist): 通用任务模型(邀N人助力得券/积分),奖励挂
--    user_coupon/point_log,记录挂本域
-- 防刷提示: 砍价/助力是被刷重灾区——帮砍/助力入口须挂风控
--    (rule_type 2高频/3异常领券/4套利特征),IP/设备/新用户限制
-- ============================================================

-- ② 推广员等级预留
ALTER TABLE `distribution_user`
  ADD COLUMN `level` TINYINT NOT NULL DEFAULT 1 COMMENT '推广员等级(V1单一等级=1;推广员>500人时启用多级比例,届时commission_rule加等级维度)' AFTER `status`;

-- ③ 砍价活动(SKU级,对齐 group_buy_item/flash_sale_item 模式)
CREATE TABLE `bargain_activity` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '砍价活动ID',
  `name`       VARCHAR(64)     NOT NULL COMMENT '活动名称',
  `spu_id`     BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID',
  `start_time` DATETIME        NOT NULL COMMENT '开始时间(含)',
  `end_time`   DATETIME        NOT NULL COMMENT '结束时间(不含)',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_spu_status` (`spu_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='砍价活动表(时段/状态;价格区间在bargain_item按SKU配置)';

CREATE TABLE `bargain_item` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '砍价场次商品ID',
  `activity_id`    BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `sku_id`         BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `original_price` DECIMAL(10,2)   NOT NULL COMMENT '起始价(=发起时售价口径)',
  `floor_price`    DECIMAL(10,2)   NOT NULL COMMENT '底价(砍到底价即可下单;floor<=original应用校验)',
  `max_cut_count`  INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '最大帮砍刀数(0=不限,金额收敛到底价)',
  `config`         JSON            NULL COMMENT '玩法参数(随机递减算法边界/首刀上限等,应用层约定)',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_sku` (`activity_id`, `sku_id`),
  KEY `idx_sku` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='砍价场次商品表(SKU级价格区间)';

CREATE TABLE `bargain_record` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '砍价单ID',
  `bargain_no`   VARCHAR(32)     NOT NULL COMMENT '砍价单号(全局唯一)',
  `item_id`      BIGINT UNSIGNED NOT NULL COMMENT '砍价场次商品ID',
  `user_id`      BIGINT UNSIGNED NOT NULL COMMENT '发起人用户ID',
  `current_price` DECIMAL(10,2)  NOT NULL COMMENT '当前价(逐刀递减,>=floor_price由应用+条件更新保证)',
  `cut_count`    INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '已砍刀数',
  `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1砍价中 2到底价待下单 3已下单 4超时失败 5已取消',
  `expire_time`  DATETIME        NOT NULL COMMENT '砍价截止时间(超时扫描)',
  `success_time` DATETIME        NULL DEFAULT NULL COMMENT '到底价时间',
  `order_no`     VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '成交订单号(下单后回填,防重复成交)',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bargain_no` (`bargain_no`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_status` (`user_id`, `status`),
  KEY `idx_status_expire` (`status`, `expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='砍价单表(当前价条件更新防并发超砍;超时/取消不可下单)';

CREATE TABLE `bargain_helper` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '帮砍记录ID',
  `record_id`       BIGINT UNSIGNED NOT NULL COMMENT '砍价单ID',
  `helper_user_id`  BIGINT UNSIGNED NOT NULL COMMENT '帮砍人用户ID',
  `cut_amount`      DECIMAL(10,2)   NOT NULL COMMENT '本刀砍掉金额(与record条件更新同事务)',
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '帮砍时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_record_helper` (`record_id`, `helper_user_id`),
  KEY `idx_helper` (`helper_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帮砍记录表(一人一刀;只追加;防刷挂风控)';

-- 订单关联砍价单(可空,普通订单零影响;砍价成交价经订单项price快照)
ALTER TABLE `trade_order`
  ADD COLUMN `bargain_record_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '砍价单ID(NULL=非砍价订单)' AFTER `group_buy_team_id`,
  ADD KEY `idx_bargain_record` (`bargain_record_id`);

-- ④ 助力(通用任务模型: 邀N人助力得券/积分)
CREATE TABLE `assist_activity` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '助力活动ID',
  `name`           VARCHAR(64)     NOT NULL COMMENT '活动名称',
  `reward_type`    TINYINT         NOT NULL COMMENT '奖励类型:1优惠券 2积分',
  `reward_ref`     BIGINT UNSIGNED NOT NULL COMMENT '奖励载体ID(如coupon.id;积分为点数存config)',
  `required_count` INT UNSIGNED    NOT NULL COMMENT '所需助力人数',
  `per_limit`      INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '每人可发起次数',
  `start_time`     DATETIME        NOT NULL COMMENT '开始时间(含)',
  `end_time`       DATETIME        NOT NULL COMMENT '结束时间(不含)',
  `config`         JSON            NULL COMMENT '扩展参数(积分数/阶梯奖励等)',
  `status`         TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`        TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='助力活动表(通用任务:邀N人得奖励)';

CREATE TABLE `assist_record` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '助力参与记录ID',
  `activity_id`  BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `user_id`      BIGINT UNSIGNED NOT NULL COMMENT '发起人用户ID',
  `helper_count` INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '已助力人数(达required_count触发发奖)',
  `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1进行中 2已完成发奖 3已过期',
  `finish_time`  DATETIME        NULL DEFAULT NULL COMMENT '完成时间',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_user` (`activity_id`, `user_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='助力参与记录表(发奖幂等:完成态不重复发)';

CREATE TABLE `assist_helper` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '助力记录ID',
  `record_id`      BIGINT UNSIGNED NOT NULL COMMENT '参与记录ID',
  `helper_user_id` BIGINT UNSIGNED NOT NULL COMMENT '助力人用户ID',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '助力时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_record_helper` (`record_id`, `helper_user_id`),
  KEY `idx_helper` (`helper_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='助力人记录表(一人一助力;只追加;防刷挂风控)';
