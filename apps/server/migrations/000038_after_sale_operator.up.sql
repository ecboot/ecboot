-- 000038_after_sale_operator: 售后单的操作人与退款失败原因（013-after-sale / 批次07）
-- 背景（research D3）: IAfterSaleLogic 的 Approve/Reject/ConfirmReceipt/RetryRefund 四个方法
--   都接收 operator 参数, 而 after_sale_order **没有任何列能承载它** → 不补就是"审计信息被静默丢弃",
--   与批次 06 评审在 AdminCancel 抓到的 Important 缺陷同型。
--   注: admin_operation_log 是请求级审计表, 但至今无人写入（其列表端点属批次 12 风控审计范围）,
--   不能作为本批载体; 新建 after_sale_log 表对本批"只需最后操作人"而言过重（宪法 V 简单优先）。
-- fail_reason: 渠道退款失败时**保留原因**并停留在"待退款", 使后台"重试退款"这条路径有据可查（FR-012）。
ALTER TABLE `after_sale_order`
  ADD COLUMN `operator_id` VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '最后操作人标识:admin:{id}(审核/确认收货/退款重试)' AFTER `reject_reason`,
  ADD COLUMN `fail_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '渠道退款失败原因(重试前保留,成功后清空)' AFTER `operator_id`;
