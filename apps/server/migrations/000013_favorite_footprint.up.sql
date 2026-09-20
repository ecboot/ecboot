-- ============================================================
-- V13 收藏与足迹域 (favorite & footprint)
-- 表: user_favorite(收藏), user_footprint(浏览足迹)
-- 均以 (user_id, spu_id) 唯一; 足迹重复浏览=更新 last_view_at
-- 足迹不软删(物理清理,保留期90天,定期任务扫描 idx(last_view_at))
-- ============================================================

CREATE TABLE `user_favorite` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `spu_id`     BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID(实时联查现价与可售状态,不快照)',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_spu` (`user_id`, `spu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户收藏表(重复收藏被唯一约束拒绝)';

CREATE TABLE `user_footprint` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '足迹ID',
  `user_id`      BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `spu_id`       BIGINT UNSIGNED NOT NULL COMMENT 'SPU ID',
  `view_count`   INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '累计浏览次数(重复浏览累加)',
  `last_view_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最近浏览时间(重复浏览更新此列,清理任务扫描)',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '首次浏览时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_spu` (`user_id`, `spu_id`),
  KEY `idx_last_view` (`last_view_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户浏览足迹表(保留期90天,物理清理,无deleted)';
