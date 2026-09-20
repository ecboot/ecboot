-- ============================================================
-- 000031 门店域 (产品定档 2026-09-20: B2C + 多门店社交电商)
-- 定位: 纯 B2C——平台统一经营商品与交易, 门店是线下载体
--   (自提/到店核销/附近门店), 无商户/商家/租户概念
--   (B2B2C 如未来需要, 另立独立项目)
-- 区划码对齐 V23 收货地址口径(GB/T 2260 六位)
-- ============================================================

CREATE TABLE `store` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '门店ID',
  `store_no`       VARCHAR(32)     NOT NULL COMMENT '门店编码(全局唯一)',
  `name`           VARCHAR(64)     NOT NULL COMMENT '门店名称',
  `province_code`  CHAR(6)         NOT NULL DEFAULT '' COMMENT '省级行政区划代码(GB/T 2260)',
  `city_code`      CHAR(6)         NOT NULL DEFAULT '' COMMENT '市级行政区划代码',
  `district_code`  CHAR(6)         NOT NULL DEFAULT '' COMMENT '区县级行政区划代码',
  `detail_address` VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '详细地址(街道门牌)',
  `longitude`      DECIMAL(10,6)   NULL DEFAULT NULL COMMENT '经度(附近门店检索;GCJ-02坐标系)',
  `latitude`       DECIMAL(10,6)   NULL DEFAULT NULL COMMENT '纬度(同上)',
  `business_hours` VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '营业时间(如 09:00-21:00)',
  `contact_phone`  VARCHAR(20)     NOT NULL DEFAULT '' COMMENT '门店电话',
  `pickup_enabled` TINYINT         NOT NULL DEFAULT 0 COMMENT '自提开关:0否 1是(下单自提点候选)',
  `sort`           INT             NOT NULL DEFAULT 0 COMMENT '排序,越小越靠前',
  `status`         TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1营业 2歇业(歇业=展示但不可自提/核销)',
  `deleted`        TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_store_no` (`store_no`),
  KEY `idx_district_status` (`district_code`, `status`),
  KEY `idx_location` (`longitude`, `latitude`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='线下门店表(B2C多门店:自提/核销/附近门店;无商户商家概念)';
