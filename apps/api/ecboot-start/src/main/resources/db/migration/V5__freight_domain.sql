-- ============================================================
-- V5 运费域 (freight domain)
-- 表: freight_template(运费模板), freight_rule(计费规则)
-- 计费方式: 按件数 / 按重量(sku.weight,无重量按1件折算)
-- 区域差异化: 按收货省匹配规则,空 regions 为全国兜底
-- 满额包邮: template.free_threshold (NULL=不包邮)
-- 商品关联: product_spu.freight_template_id (NULL=包邮)
-- 下单时计算运费并快照进 trade_order (freight_amount+template_id)
-- ============================================================

CREATE TABLE `freight_template` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name`            VARCHAR(64)     NOT NULL COMMENT '模板名称(如:华东包邮-全国2件10元)',
  `charge_type`     TINYINT         NOT NULL DEFAULT 1 COMMENT '计费方式:1按件数 2按重量',
  `free_threshold`  DECIMAL(10,2)   NULL DEFAULT NULL COMMENT '满额包邮阈值(订单商品金额≥此值免运费,NULL=不包邮)',
  `status`          TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用(停用后新订单不可用,已下单不受影响)',
  `deleted`         TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='运费模板表(SPU关联;每模板必含一条空regions兜底规则,应用层保证)';

CREATE TABLE `freight_rule` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规则ID',
  `template_id`    BIGINT UNSIGNED NOT NULL COMMENT '所属模板ID',
  `region_codes`   JSON            NOT NULL COMMENT '适用省级行政区代码列表(如["310000","330000"];空数组=全国兜底规则,每模板至多一条)',
  `first_unit`     INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '首段额度:按件=首件数;按重=首重克数',
  `first_fee`      DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '首段运费',
  `continue_unit`  INT UNSIGNED    NOT NULL DEFAULT 1 COMMENT '续段步长:按件=每续N件;按重=每续N克',
  `continue_fee`   DECIMAL(10,2)   NOT NULL DEFAULT 0.00 COMMENT '续段单价(每步长追加费用)',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_template` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='运费计费规则表(区域差异化)';

-- 运费计算规则(应用层实现):
--   1. 订单商品金额 ≥ template.free_threshold → 运费0(满额包邮)
--   2. 按 charge_type 折算计费量 Q: 按件=订单总件数; 按重=Σ(sku.weight*qty),weight缺失的SKU按1件=1单位折算
--   3. 匹配收货省代码对应的 rule(无匹配则空regions兜底规则)
--   4. 运费 = first_fee + Q≤first_unit?0:CEIL((Q-first_unit)/continue_unit)*continue_fee
