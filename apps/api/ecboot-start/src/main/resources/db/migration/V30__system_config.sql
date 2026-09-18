-- ============================================================
-- V30 系统配置表 (产品决策 2026-09-18: 推广员等级触发条件可配置)
-- 动机: 散落在列注释里的运营参数(等级阈值/归因窗口/订单超时/
--   自动收货/结算保护期)全部硬编码——收敛建配置表,运营可调
-- 语义: status=0停用时应用回退代码内置默认值(配置是覆盖层,
--   不是唯一真源——配置表异常不应导致业务不可用)
-- ============================================================

CREATE TABLE `system_config` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `code`        VARCHAR(64)     NOT NULL COMMENT '配置编码(点分命名空间,如 distribution.level.threshold)',
  `value`       VARCHAR(1024)   NOT NULL COMMENT '配置值(标量或JSON对象)',
  `value_type`  TINYINT         NOT NULL DEFAULT 1 COMMENT '值类型:1整数 2小数 3字符串 4布尔 5JSON对象',
  `name`        VARCHAR(64)     NOT NULL COMMENT '配置名称(后台展示)',
  `description` VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '配置说明(含代码默认值,便于回退排查)',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用(停用回退应用代码内置默认值)',
  `deleted`     TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统配置表(运营参数配置化;配置=覆盖层,停用/缺失回退代码默认值)';

-- 种子配置(初始值=设计文档既有默认值;调整改行,不改代码)
INSERT INTO `system_config` (`code`, `value`, `value_type`, `name`, `description`) VALUES
('distribution.level.threshold',   '500', 1, '推广员等级体系启用阈值', '通过审核推广员人数超过该值时启用多级比例(届时commission_rule加等级维度);代码默认500'),
('attribution.window.days',        '7',   1, '分享归因窗口(天)',       '分享后N天内下单归因给分享人;代码默认7'),
('order.timeout.minutes',          '30',  1, '待付款订单超时(分钟)',   '超时自动取消并释放库存;代码默认30'),
('order.auto_confirm.days',        '7',   1, '自动确认收货(天)',       '发货后N天自动完成订单;代码默认7'),
('commission.settle.protect_days', '7',   1, '佣金结算保护期(天)',      '确认收货后N天保护期满结算,保护期退款不结算;代码默认7'),
('account_payment.enabled',        'true', 4, '佣金余额消费开关',       '是否允许下单使用佣金余额抵扣;代码默认true');

-- 等级列注释解除硬编码阈值,指向配置
ALTER TABLE `distribution_user`
  MODIFY COLUMN `level` TINYINT NOT NULL DEFAULT 1 COMMENT '推广员等级(V1单一等级=1;多级启用阈值见system_config: distribution.level.threshold)';
