-- ============================================================
-- V23 业务评审修复 (business review fixes, PM+技术双视角评审)
-- P0-① 地址区划代码：运费规则 region_codes 匹配断裂修复
-- P0-② 售后退货数量：买3退1的退款/回补/冲销依据
-- P0-③ 秒杀活动关联：限购校验与活动对账的数据落点
-- P1-④ 满额包邮区域例外 / P1-⑤ 关系链锁定 / P1-⑦ 规格组合结构唯一
-- 附: point_log 业务类型枚举扩位(纯注释)
-- ============================================================

-- P0-① 收货地址增加国标区划代码(名称列保留=展示快照;代码=运费匹配唯一口径)
ALTER TABLE `user_address`
  ADD COLUMN `province_code` CHAR(6) NOT NULL DEFAULT '' COMMENT '省级行政区划代码(GB/T 2260六位,运费规则匹配口径;空=历史数据待补)' AFTER `province`,
  ADD COLUMN `city_code`     CHAR(6) NOT NULL DEFAULT '' COMMENT '市级行政区划代码' AFTER `city`,
  ADD COLUMN `district_code` CHAR(6) NOT NULL DEFAULT '' COMMENT '区县级行政区划代码' AFTER `district`;

-- P0-② 售后单退货数量(退款金额摊算/库存回补/佣金冲销的数量依据)
ALTER TABLE `after_sale_order`
  ADD COLUMN `quantity` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '退货/退款数量(≤订单项quantity,应用校验;部分退货的金额与回补依据)' AFTER `currency`;

-- P0-③ 订单项秒杀活动关联(限购校验:SUM(quantity)按用户+场次;活动对账)
ALTER TABLE `trade_order_item`
  ADD COLUMN `flash_sale_item_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '秒杀场次商品ID(NULL=非秒杀;限购校验与秒杀订单识别依据)' AFTER `sku_specs`,
  ADD KEY `idx_flash_item` (`flash_sale_item_id`);

-- P1-④ 满额包邮区域例外(偏远省份不参与包邮)
ALTER TABLE `freight_template`
  ADD COLUMN `free_exclude_codes` JSON NULL COMMENT '不参与满额包邮的省级区划代码列表(如["650000","540000"]);NULL=全国均可包邮' AFTER `free_threshold`;

-- P1-⑤ 关系链锁定(绑定窗口规则:保护期内可换绑,期满锁定;换绑=UPDATE本行,UNIQUE(user_id)不受影响)
ALTER TABLE `user_relation`
  ADD COLUMN `locked`    TINYINT  NOT NULL DEFAULT 0 COMMENT '锁定:0保护期内(可换绑) 1已锁定(绑定窗口期满,防抢人纠纷)' AFTER `bind_time`,
  ADD COLUMN `lock_time` DATETIME NULL DEFAULT NULL COMMENT '锁定时间' AFTER `locked`;

-- P1-⑦ SKU规格组合唯一结构强制(生成列哈希;依赖应用层JSON序列化键序稳定,键序不一致时退化为不误伤的弱约束)
ALTER TABLE `product_sku`
  ADD COLUMN `specs_hash` CHAR(64) AS (SHA2(`specs`, 256)) STORED COMMENT '规格组合哈希(同SPU内组合唯一;应用序列化稳定时生效)' AFTER `specs`,
  ADD UNIQUE KEY `uk_spu_specs_hash` (`spu_id`, `specs_hash`);

-- 附: point_log 业务类型枚举扩位(纯注释,不改数据)
ALTER TABLE `point_log`
  MODIFY COLUMN `biz_type` TINYINT NOT NULL COMMENT '业务类型:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励';
