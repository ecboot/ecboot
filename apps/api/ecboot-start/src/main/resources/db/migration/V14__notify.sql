-- ============================================================
-- V14 通知域 (notification)
-- 表: notify_task(异步通知任务,一渠道一行 research D4),
--     user_message(站内消息,必达兜底)
-- 渠道被用户关闭→不建任务行(站内信例外必达)
-- 发送失败自动重试(有上限),不影响业务主流程(spec FR-012)
-- ============================================================

CREATE TABLE `notify_task` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `user_id`         BIGINT UNSIGNED NOT NULL COMMENT '接收用户ID',
  `channel`         TINYINT         NOT NULL COMMENT '渠道:1小程序订阅消息 2短信 3站内信',
  `biz_type`        TINYINT         NOT NULL COMMENT '业务类型:1订单 2营销 3售后',
  `biz_no`          VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '关联业务单号(如订单号)',
  `template_code`   VARCHAR(64)     NOT NULL COMMENT '消息模板编码',
  `params`          JSON            NULL COMMENT '模板参数(键值对)',
  `status`          TINYINT         NOT NULL DEFAULT 10 COMMENT '状态:10待发送 20已发送 30失败 40达重试上限 50已跳过(渠道关闭)',
  `retry_count`     INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '已重试次数',
  `next_retry_time` DATETIME        NULL DEFAULT NULL COMMENT '下次重试时间(失败退避后回填)',
  `sent_time`       DATETIME        NULL DEFAULT NULL COMMENT '发送成功时间',
  `fail_reason`     VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '最后失败原因',
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`, `created_at`),
  KEY `idx_status_retry` (`status`, `next_retry_time`),
  KEY `idx_biz_no` (`biz_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='异步通知任务表(一渠道一行,重试有上限,终态40/50)';

CREATE TABLE `user_message` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '接收用户ID',
  `title`       VARCHAR(128)    NOT NULL COMMENT '标题',
  `content`     VARCHAR(2048)   NOT NULL COMMENT '内容',
  `biz_type`    TINYINT         NOT NULL DEFAULT 1 COMMENT '业务类型:1订单 2营销 3售后',
  `biz_no`      VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '关联业务单号(点击跳转)',
  `is_read`     TINYINT         NOT NULL DEFAULT 0 COMMENT '已读:0否 1是',
  `read_time`   DATETIME        NULL DEFAULT NULL COMMENT '阅读时间',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_read_time` (`user_id`, `is_read`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='站内消息表(必达兜底渠道,不支持删除仅已读)';
