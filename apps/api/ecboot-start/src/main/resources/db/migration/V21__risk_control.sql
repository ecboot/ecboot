-- ============================================================
-- V21 风控域 (risk control)
-- 表: risk_rule(规则), risk_record(事件,只追加)
-- 注: 站内搜索不建表(spec显式决定——商品既有结构足够,
--     检索能力属技术实现层,契约边界封口)
-- 事件多态关联: object_type+object_no 定位订单/券/提现/售后
-- 申诉通道: appeal_status 状态机,后台白名单可解除拦截
-- ============================================================

CREATE TABLE `risk_rule` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规则ID',
  `name`          VARCHAR(64)     NOT NULL COMMENT '规则名称(如 黑名单/高频下单/佣金套利特征)',
  `rule_type`     TINYINT         NOT NULL COMMENT '规则类型:1黑名单 2高频下单 3异常领券 4佣金套利特征',
  `condition_expr` VARCHAR(512)   NOT NULL COMMENT '规则条件描述(如 1分钟内下单>10次)',
  `action`        TINYINT         NOT NULL DEFAULT 1 COMMENT '处置:1拦截 2标记(放行但留痕)',
  `status`        TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`       TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='风控规则表(营销与分销资金面的拦截规则)';

CREATE TABLE `risk_record` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '事件ID',
  `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '命中用户ID',
  `rule_id`       BIGINT UNSIGNED NOT NULL COMMENT '命中规则ID',
  `object_type`   TINYINT         NOT NULL COMMENT '关联对象类型:1订单 2优惠券 3提现 4售后',
  `object_no`     VARCHAR(32)     NOT NULL COMMENT '关联对象单号(订单号/券ID/提现单号等)',
  `action`        TINYINT         NOT NULL DEFAULT 1 COMMENT '处置结果:1拦截 2标记',
  `appeal_status` TINYINT         NOT NULL DEFAULT 0 COMMENT '申诉状态:0无 1申诉中 2申诉通过(解除拦截) 3申诉驳回',
  `remark`        VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '备注(命中明细/申诉结论)',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '命中时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`, `created_at`),
  KEY `idx_object` (`object_type`, `object_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='风控事件表(只追加,拦截可追溯,申诉四态)';
