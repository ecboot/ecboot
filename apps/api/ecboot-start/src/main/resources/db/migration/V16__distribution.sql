-- ============================================================
-- V16 分销域 (distribution) — 8 表
-- 合规红线: 关系链两级封顶(禁止传销条例), 仅 inviter_id 单列,
--           结构上无法表达三级 → 详见 docs/adr/0003
-- 佣金基数 = 订单项实付(pay_amount,不含运费);命中 商品>分类
-- 账户余额可负(欠款抵扣); 提现渠道单号唯一(幂等,同支付总则)
-- ============================================================

-- 用户关系链(一人一条,两级封顶的结构载体)
CREATE TABLE `user_relation` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系ID',
  `user_id`      BIGINT UNSIGNED NOT NULL COMMENT '用户ID(一人至多一条关系)',
  `inviter_id`   BIGINT UNSIGNED NOT NULL COMMENT '直接上级(邀请人);二级=上级的上级,两次单列查询解析;无祖父列/路径列(ADR-0003)',
  `bind_channel` TINYINT         NOT NULL DEFAULT 1 COMMENT '绑定方式:1分享链接 2邀请码',
  `bind_time`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user` (`user_id`),
  KEY `idx_inviter` (`inviter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户邀请关系链(仅两级,法规红线结构强制)';

-- 推广员资质
CREATE TABLE `distribution_user` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '推广员ID',
  `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1待审核 2通过 3冻结(冻结不产生新佣金,存量可提现)',
  `apply_time`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `audit_time`  DATETIME        NULL DEFAULT NULL COMMENT '审核时间',
  `deleted`     TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='推广员资质表(后台审核生效)';

-- 佣金规则(分类默认+商品覆盖,命中优先级 商品>分类,未命中不计佣)
CREATE TABLE `commission_rule` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规则ID',
  `scope_type`   TINYINT         NOT NULL COMMENT '作用域:1分类(默认) 2商品(覆盖)',
  `scope_id`     BIGINT UNSIGNED NOT NULL COMMENT '作用域目标ID(分类ID或SPU ID)',
  `level1_rate`  DECIMAL(5,2)    NOT NULL DEFAULT 0.00 COMMENT '一级佣金比例%(直接邀请人,0-100)',
  `level2_rate`  DECIMAL(5,2)    NOT NULL DEFAULT 0.00 COMMENT '二级佣金比例%(邀请人的邀请人,0-100)',
  `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`      TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_scope` (`scope_type`, `scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='佣金规则表(商品>分类命中,基数=订单项实付)';

-- 佣金记录(订单项粒度,含冲销)
CREATE TABLE `commission_record` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '佣金记录ID',
  `order_no`            VARCHAR(32)     NOT NULL COMMENT '订单号',
  `order_item_id`       BIGINT UNSIGNED NOT NULL COMMENT '订单项ID',
  `beneficiary_user_id` BIGINT UNSIGNED NOT NULL COMMENT '受益人用户ID',
  `level`               TINYINT         NOT NULL COMMENT '层级:1直接邀请 2间接邀请',
  `base_amount`         DECIMAL(10,2)   NOT NULL COMMENT '计佣基数(=订单项实付pay_amount)',
  `rate`                DECIMAL(5,2)    NOT NULL COMMENT '命中比例%',
  `amount`              DECIMAL(10,2)   NOT NULL COMMENT '佣金金额(冲销记录为负值)',
  `status`              TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1待结算 2已结算 3已失效(保护期退款) 4欠款冲销中(结算后退款)',
  `settle_time`         DATETIME        NULL DEFAULT NULL COMMENT '结算时间(确认收货+7天保护期满)',
  `reversal_record_id`  BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '冲销关联:本记录被哪条负额冲销记录回指(原记录侧)',
  `reversal_of_id`      BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '冲销关联:本冲销记录冲销的原记录ID(冲销记录侧)',
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order` (`order_no`),
  KEY `idx_user_status` (`beneficiary_user_id`, `status`),
  KEY `idx_item` (`order_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='佣金记录表(待结算→已结算;退款冲销,负额关联原记录)';

-- 用户佣金账户(余额可负=欠款抵扣)
CREATE TABLE `user_account` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '账户ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `balance`    DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '可用余额(有符号,欠款为负,后续入账抵扣)',
  `frozen`     DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '冻结金额(提现审核/打款中)',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户佣金账户表(可负余额,禁止UNSIGNED)';

-- 账户流水(只追加,双向可追溯)
CREATE TABLE `account_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `biz_type`      TINYINT         NOT NULL COMMENT '业务类型:1佣金入账 2提现冻结 3提现完成 4提现失败回退 5冲销扣回',
  `amount`        DECIMAL(10,2)   NOT NULL COMMENT '变动金额(有符号:入账正,冻结/扣回负)',
  `balance_after` DECIMAL(10,2)   NOT NULL COMMENT '变动后余额快照(对账用)',
  `frozen_after`  DECIMAL(10,2)   NOT NULL COMMENT '变动后冻结快照(对账用)',
  `biz_no`        VARCHAR(32)     NOT NULL COMMENT '关联业务单号(commission_record.id或withdraw_order.withdraw_no)',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`, `id`),
  KEY `idx_biz_no` (`biz_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='佣金账户流水表(只追加,余额/冻结双快照)';

-- 提现单(渠道单号唯一幂等,同支付回调总则)
CREATE TABLE `withdraw_order` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '提现单ID',
  `withdraw_no`      VARCHAR(32)     NOT NULL COMMENT '提现单号(全局唯一)',
  `user_id`          BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `amount`           DECIMAL(10,2)   NOT NULL COMMENT '提现金额',
  `withdraw_channel` TINYINT         NOT NULL DEFAULT 1 COMMENT '渠道:1微信商家转账(V1单渠道)',
  `channel_order_no` VARCHAR(64)     NULL DEFAULT NULL COMMENT '渠道打款单号(回填;与渠道联合唯一,防重复打款)',
  `status`           TINYINT         NOT NULL DEFAULT 10 COMMENT '状态:10待审核 20审核通过 30打款中 40成功 50审核拒绝 60打款失败已回退',
  `audit_time`       DATETIME        NULL DEFAULT NULL COMMENT '审核时间',
  `pay_time`         DATETIME        NULL DEFAULT NULL COMMENT '打款成功时间',
  `fail_reason`      VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '拒绝/失败原因',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_withdraw_no` (`withdraw_no`),
  UNIQUE KEY `uk_channel_order` (`withdraw_channel`, `channel_order_no`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='佣金提现单表(打款回调条件更新WHERE status=30幂等)';

-- 邀请注册激励记录(一新用户仅一次)
CREATE TABLE `invite_record` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '激励记录ID',
  `new_user_id`  BIGINT UNSIGNED NOT NULL COMMENT '新用户ID(仅可被激励一次)',
  `inviter_id`   BIGINT UNSIGNED NOT NULL COMMENT '邀请人用户ID',
  `reward_type`  TINYINT         NOT NULL DEFAULT 1 COMMENT '奖励类型:1优惠券',
  `reward_ref`   BIGINT UNSIGNED NOT NULL COMMENT '奖励载体ID(如user_coupon.id)',
  `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1已发放(发放与订单同事务,幂等)',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_new_user` (`new_user_id`),
  KEY `idx_inviter` (`inviter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='邀请注册激励记录表(一人一次,发放幂等)';
