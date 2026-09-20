-- ============================================================
-- V20 满减活动域 (full reduction)
-- 表: promotion_activity, promotion_activity_ladder,
--     promotion_activity_scope + 订单优惠三构成列
-- 范围表达(research D3): 关系表 scope_type 1全场/2分类/3商品,
--   全场=target_id NULL 单行; idx(scope_type,target_id)支撑下单
--   热路径反查"哪些活动覆盖此商品/分类"
-- 优惠三构成恒等式(research D8):
--   trade_order.promotion_amount ≡ coupon_amount
--     + full_reduction_amount + point_amount
--   (订单项分摊合计=头合计,尾差记末行;校验SQL见quickstart场景三)
-- 档位命中: 自动选用户可享最优档;与券叠加: 先满减后用券
-- ============================================================

CREATE TABLE `promotion_activity` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '活动ID',
  `name`       VARCHAR(64)     NOT NULL COMMENT '活动名称(如 全场满减)',
  `start_time` DATETIME        NOT NULL COMMENT '开始时间(含)',
  `end_time`   DATETIME        NOT NULL COMMENT '结束时间(不含)',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status_time` (`status`, `start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='满减活动表(时段/状态;范围由scope表表达)';

CREATE TABLE `promotion_activity_ladder` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '档位ID',
  `activity_id`      BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `threshold_amount` DECIMAL(10,2)   NOT NULL COMMENT '门槛金额(满X元,范围内商品合计)',
  `discount_amount`  DECIMAL(10,2)   NOT NULL COMMENT '优惠金额(减Y元)',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_threshold` (`activity_id`, `threshold_amount`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='满减档位表(同活动同门槛唯一,多档如满100-10/满200-30)';

CREATE TABLE `promotion_activity_scope` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '范围ID',
  `activity_id` BIGINT UNSIGNED NOT NULL COMMENT '活动ID',
  `scope_type`  TINYINT         NOT NULL COMMENT '范围类型:1全场 2分类 3商品',
  `target_id`   BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '目标ID(scope_type=2分类ID/3商品SPU_ID;全场为NULL,每活动至多一条全场行)',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_activity_scope` (`activity_id`, `scope_type`, `target_id`),
  KEY `idx_target` (`scope_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='满减活动范围表(下单热路径反查命中,JSON不可索引故拆表)';

-- 订单层优惠三构成(合计列 promotion_amount 保留,恒等于三明细之和)
ALTER TABLE `trade_order`
  ADD COLUMN `coupon_amount`         DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '优惠券抵扣金额(promotion_amount构成之一)' AFTER `point_used`,
  ADD COLUMN `full_reduction_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '满减优惠金额(promotion_amount构成之一)' AFTER `coupon_amount`,
  ADD COLUMN `promotion_activity_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '命中的满减活动ID(NULL=未参与)' AFTER `full_reduction_amount`;

ALTER TABLE `trade_order_item`
  ADD COLUMN `coupon_amount`         DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '优惠券抵扣行分摊' AFTER `point_amount`,
  ADD COLUMN `full_reduction_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '满减优惠行分摊(两项合计=订单头对应明细,尾差记末行)' AFTER `coupon_amount`;
