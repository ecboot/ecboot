-- ============================================================
-- V31 店铺/商家/商户域 (产品决策 2026-09-18: 预留式多商家)
-- 统一语言(术语三分,见 CONTEXT.md):
--   merchant(商户)=结算与资质主体——管钱(支付商户号/进件/分账)
--   seller(商家)  =交易卖方——管交易(商品/订单/售后归属方)
--   shop(店铺)    =C端门面——管门面(装修/公告/客服入口)
-- 关系: 1—1—1 三视角(V1 自营单主体,各一行种子);
--   维度列一律指向 seller(交易归属);多商家开放时关系可放宽
-- 分销联动: 商家商品与自营商品同规则参与分享归因与两级佣金
-- B形态演进清单(届时再建,不阻塞): 入驻审核流/保证金/
--   seller_user商家端账号/平台-商家分账结算/购物车按商家分组
--   拆单/店铺级装修配置/商家自营销
-- ============================================================

CREATE TABLE `merchant` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商户ID(1=自营,种子行)',
  `merchant_no`   VARCHAR(32)     NOT NULL COMMENT '商户编码(全局唯一)',
  `name`          VARCHAR(64)     NOT NULL COMMENT '商户名称(企业主体)',
  `wx_mch_id`     VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '微信支付商户号(结算主体标识;平台分账接入时启用,自营店共用平台商户号则为空)',
  `status`        TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1正常 2冻结 3待审核(入驻预留)',
  `contact_name`  VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '联系人',
  `contact_phone` VARCHAR(20)     NOT NULL DEFAULT '' COMMENT '联系电话',
  `deleted`       TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_no` (`merchant_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商户表(结算与资质主体,管钱;V1自营单主体,1:1:1三视角)';

CREATE TABLE `seller` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商家ID(1=自营,种子行;商品/订单/售后的seller_id指向此表)',
  `seller_no`   VARCHAR(32)     NOT NULL COMMENT '商家编码(全局唯一)',
  `merchant_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属商户ID(1:1,V1自营;关系可放宽预留)',
  `name`        VARCHAR(64)     NOT NULL COMMENT '商家名称(运营主体)',
  `type`        TINYINT         NOT NULL DEFAULT 1 COMMENT '类型:1自营 2入驻商家(V1仅自营,入驻态预留)',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1正常 2冻结 3待审核(入驻预留)',
  `deleted`     TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_seller_no` (`seller_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商家表(交易卖方,管交易;V1自营单主体)';

CREATE TABLE `shop` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '店铺ID(1=自营旗舰店,种子行)',
  `shop_no`              VARCHAR(32)     NOT NULL COMMENT '店铺编码(全局唯一)',
  `seller_id`            BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属商家ID(1:1,V1自营)',
  `name`                 VARCHAR(64)     NOT NULL COMMENT '店铺名称(对外展示)',
  `logo`                 VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '店铺Logo URL',
  `description`          VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '店铺简介',
  `customer_service_url` VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '客服入口(URL/企微二维码)',
  `status`               TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1营业 2歇业(歇业=可浏览不可下单,应用层控制)',
  `deleted`              TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`           DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`           DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shop_no` (`shop_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='店铺表(C端门面,管门面;V1自营单主体)';

-- 种子: 自营三视角各一行(merchant/seller/shop id=1)
INSERT INTO `merchant` (`id`, `merchant_no`, `name`) VALUES (1, 'MSELF0001', 'ECBOOT自营商户');
INSERT INTO `seller` (`id`, `seller_no`, `merchant_id`, `name`, `type`) VALUES (1, 'SELF0001', 1, 'ECBOOT自营商家', 1);
INSERT INTO `shop` (`id`, `shop_no`, `seller_id`, `name`, `description`) VALUES (1, 'SHOP0001', 1, 'ECBOOT自营旗舰店', '平台自营店铺');

-- seller 维度注入(默认1=自营,既有流程零感知;覆盖 商品/订单/售后/运费 四个归属核心)
ALTER TABLE `product_spu`
  ADD COLUMN `seller_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属商家(1=自营;多商家预留)' AFTER `brand_id`,
  ADD KEY `idx_seller` (`seller_id`);

ALTER TABLE `trade_order`
  ADD COLUMN `seller_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '订单归属商家(1=自营;V1不拆单,多商家上线时购物车按商家分组结算)' AFTER `user_id`,
  ADD KEY `idx_seller_time` (`seller_id`, `created_at`);

ALTER TABLE `after_sale_order`
  ADD COLUMN `seller_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '售后处理商家(1=自营;商家后台按此筛选)' AFTER `user_id`,
  ADD KEY `idx_seller` (`seller_id`);

ALTER TABLE `freight_template`
  ADD COLUMN `seller_id` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '运费模板归属商家(多商家时各商家自建模板)' AFTER `name`,
  ADD KEY `idx_seller` (`seller_id`);
