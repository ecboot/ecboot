-- 000036_user_notify_preference: 会员通知偏好（011-member-center）
-- 背景: 勘察确认偏好无任何既有存储（无表无列）, 而接口语义（Enqueue 按偏好拆渠道）要求持久化。
-- 语义: 未建行 = 默认全开（不预置行）; 站内信不受偏好控制（仅小程序订阅/短信两渠道）。
CREATE TABLE `user_notify_preference` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '偏好ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '会员ID',
  `channel`    TINYINT         NOT NULL COMMENT '渠道:1小程序订阅 2短信',
  `enabled`    TINYINT         NOT NULL DEFAULT 1 COMMENT '是否接收:1接收 0关闭',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_channel` (`user_id`, `channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员通知偏好(会员×渠道唯一;未建行=默认全开)';
