-- ============================================================
-- V3 库存域 (inventory domain)
-- 表: inventory(库存), inventory_log(库存流水)
-- 模型: total(总=可售+锁定), locked(下单未支付锁定),
--       可售 = total - locked (推导值,不落库)
-- 时机: 下单锁定 / 支付核销 / 取消超时释放 / 售后回补
-- 防超卖: 原子条件更新(单行行锁), CHECK(locked<=total) 兜底
-- 热点SKU: 可启用 Redis Lua 预扣前置挡板(见 docs/schema-design.md)
-- ============================================================

CREATE TABLE `inventory` (
  `sku_id`     BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID(与product_sku 1:1,主键即关联)',
  `total`      INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '总库存=可售+已锁定',
  `locked`     INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '下单锁定(未支付),可售=total-locked',
  `warn_count` INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '库存预警阈值(低于触发后台提醒,预留)',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`sku_id`),
  CONSTRAINT `chk_locked_le_total` CHECK (`locked` <= `total`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='SKU库存表(行锁热点,变更必写inventory_log)';

CREATE TABLE `inventory_log` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '流水ID',
  `sku_id`       BIGINT UNSIGNED NOT NULL COMMENT 'SKU ID',
  `order_no`     VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '关联订单号(后台调整等无单操作为空)',
  `change_type`  TINYINT         NOT NULL COMMENT '变动类型:1下单锁定 2支付核销 3取消释放 4超时释放 5后台调整 6售后回补',
  `quantity`     INT UNSIGNED    NOT NULL COMMENT '变动数量(绝对值,方向由change_type决定)',
  `total_after`  INT UNSIGNED    NOT NULL COMMENT '变动后total快照(对账用)',
  `locked_after` INT UNSIGNED    NOT NULL COMMENT '变动后locked快照(对账用)',
  `operator`     VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '操作者:system/user:{id}/admin:{id}',
  `remark`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '备注',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sku` (`sku_id`, `id`),
  KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='库存变动流水(只追加,对账与审计依据)';
