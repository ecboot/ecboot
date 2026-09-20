-- ============================================================
-- V15 物流域 (logistics)
-- 表: logistics_company(物流公司字典)
-- 与 trade_order.deliver_company(存code)衔接;
-- 停用保留记录(软删/状态位),历史订单引用不受影响
-- 一单多包裹拆单为演进预留,不在本期(spec FR-015)
-- ============================================================

CREATE TABLE `logistics_company` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '物流公司ID',
  `code`          VARCHAR(32)     NOT NULL COMMENT '编码(唯一,订单deliver_company存此值)',
  `name`          VARCHAR(64)     NOT NULL COMMENT '公司名称(如 顺丰速运)',
  `tracking_rule` VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '运单号校验规则描述(如 SF+12位数字)',
  `sort`          INT             NOT NULL DEFAULT 0 COMMENT '排序',
  `status`        TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`       TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='物流公司字典表(发货选择,停用保留)';
