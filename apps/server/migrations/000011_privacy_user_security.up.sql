-- ============================================================
-- V11 隐私安全域 (privacy & user security)
-- 表: user(改造), user_login_log(新增)
-- 手机号加密(research D2/spec FR-004):
--   phone 列语义改为密文(加密存储), 检索走 phone_hash 哈希盲索引
--   仅完整手机号精确命中, 不提供模糊/尾号检索(个保法最小化)
-- ⚠️ 存量数据环境: 须先执行应用侧密文化任务(生成密文+哈希回填)
--    再执行本迁移; 当前脚手架期库内无用户数据, 无数据步骤
-- ============================================================

ALTER TABLE `user`
  MODIFY COLUMN `phone` VARCHAR(256) NOT NULL COMMENT '手机号密文(加密存储,算法由应用层确定)',
  ADD COLUMN `phone_hash` CHAR(64) NOT NULL COMMENT '手机号哈希(加盐SHA-256,登录精确检索与唯一性)' AFTER `phone`,
  DROP INDEX `uk_phone`,
  ADD UNIQUE KEY `uk_phone_hash` (`phone_hash`);

CREATE TABLE `user_login_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `login_channel` TINYINT         NOT NULL DEFAULT 1 COMMENT '登录渠道:1微信小程序 2H5',
  `login_status`  TINYINT         NOT NULL COMMENT '结果:1成功 2失败',
  `ip`            VARCHAR(45)     NOT NULL DEFAULT '' COMMENT '来源IP(IPv6最长45字符)',
  `user_agent`    VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '浏览器/客户端UA',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录日志(只追加,用户可查近30天)';
