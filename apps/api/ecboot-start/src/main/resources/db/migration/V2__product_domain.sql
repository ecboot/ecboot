-- ============================================================
-- V2 商品域 (product domain)
-- 表: product_category(分类), product_brand(品牌),
--     product_spu(SPU), product_sku(SKU)
-- 规格存储: SPU 存规格定义 spec_definitions JSON,
--          SKU 存规格组合快照 specs JSON
-- 可售 = SPU上架 AND SKU启用 AND 库存>0 (两级上下架)
-- 运费: SPU 关联 freight_template (NULL=包邮)
-- ============================================================

CREATE TABLE `product_category` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `parent_id`  BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID,0为根(固定三级:1/2/3)',
  `name`       VARCHAR(64)     NOT NULL COMMENT '分类名称',
  `icon`       VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '分类图标URL',
  `level`      TINYINT         NOT NULL COMMENT '层级:1一级 2二级 3三级',
  `sort`       INT             NOT NULL DEFAULT 0 COMMENT '同级排序,越小越靠前',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0禁用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent` (`parent_id`),
  KEY `idx_level_sort` (`level`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品分类表(自引用树,三级)';

CREATE TABLE `product_brand` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID',
  `name`        VARCHAR(64)     NOT NULL COMMENT '品牌名称(存活期间唯一,软删后仍占用)',
  `logo`        VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '品牌Logo URL',
  `description` VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '品牌简介',
  `sort`        INT             NOT NULL DEFAULT 0 COMMENT '排序,越小越靠前',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0禁用',
  `deleted`     TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品品牌表';

CREATE TABLE `product_spu` (
  `id`                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'SPU ID',
  `spu_no`             VARCHAR(32)     NOT NULL COMMENT 'SPU业务编码(全局唯一,对外展示用)',
  `name`               VARCHAR(128)    NOT NULL COMMENT '商品名称(SPU级)',
  `sub_title`          VARCHAR(256)    NOT NULL DEFAULT '' COMMENT '副标题/卖点',
  `category_id`        BIGINT UNSIGNED NOT NULL COMMENT '所属分类ID(三级分类)',
  `brand_id`           BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '品牌ID,可为空(无品牌商品)',
  `freight_template_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '运费模板ID(NULL=包邮)',
  `images`             JSON            NOT NULL COMMENT '主图+轮播图URL数组,按序存储',
  `description`        TEXT            NULL COMMENT '图文详情(富文本)',
  `video_url`          VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '主图视频URL(预留)',
  `spec_definitions`   JSON            NULL COMMENT '规格定义:[{"name":"颜色","values":["黑","白"]}]',
  `attributes`         JSON            NULL COMMENT '非销售属性键值对(材质/产地等)',
  `status`             TINYINT         NOT NULL DEFAULT 0 COMMENT '上架状态:0下架 1上架(默认下架)',
  `sale_count`         INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '累计销量(异步冗余,支付成功事件累加,排序用)',
  `deleted`            TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是(删除后历史订单不受影响,依赖订单项快照)',
  `created_at`         DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`         DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu_no` (`spu_no`),
  KEY `idx_category_status` (`category_id`, `status`),
  KEY `idx_brand` (`brand_id`),
  KEY `idx_status_sale` (`status`, `sale_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品SPU表(标准产品单元)';

CREATE TABLE `product_sku` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `sku_no`       VARCHAR(32)     NOT NULL COMMENT 'SKU业务编码(全局唯一)',
  `spu_id`       BIGINT UNSIGNED NOT NULL COMMENT '所属SPU ID',
  `name`         VARCHAR(128)    NOT NULL DEFAULT '' COMMENT 'SKU名称=SPU名+规格串(冗余生成,展示用)',
  `specs`        JSON            NOT NULL COMMENT '规格组合快照:{"颜色":"黑","尺码":"M"}(同SPU内组合唯一由应用层保证)',
  `price`        DECIMAL(10,2)   NOT NULL COMMENT '现售价(元)',
  `line_price`   DECIMAL(10,2)   NULL DEFAULT NULL COMMENT '划线价/原价(营销展示,可为空)',
  `cost_price`   DECIMAL(10,2)   NULL DEFAULT NULL COMMENT '成本价(后台可见,毛利分析预留)',
  `image`        VARCHAR(512)    NOT NULL DEFAULT '' COMMENT 'SKU图URL(空则用SPU主图)',
  `weight`       DECIMAL(10,2)   NULL DEFAULT NULL COMMENT '重量(克,运费按重计费时使用)',
  `barcode`      VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '条形码(预留)',
  `sort`         INT             NOT NULL DEFAULT 0 COMMENT '展示排序',
  `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0禁用(断码/失效单品)',
  `deleted`      TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_no` (`sku_no`),
  KEY `idx_spu` (`spu_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品SKU表(库存进出计价单元)';
