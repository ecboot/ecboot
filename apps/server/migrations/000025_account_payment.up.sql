-- ============================================================
-- V25 佣金余额可消费 (产品决策 2026-09-18, 评审遗留项2)
-- 资金语义: 余额抵扣=用户资产消耗(非营销优惠), 不改既有
--   恒等式, 只加一层: 现金实付 = pay_amount - account_amount
--   (pay_order.amount 即现金实付, 渠道单只收现金部分)
-- 既有恒等式(全部保持):
--   promotion_amount ≡ coupon + full_reduction + point (头/项两层)
--   pay_amount ≡ total_amount - promotion_amount + freight_amount
-- 新增恒等式:
--   pay_order.amount ≡ trade_order.pay_amount - trade_order.account_amount
-- 资金流程(与提现同构, 条件更新):
--   下单冻结(balance→frozen) → 支付成功核销(frozen扣减) →
--   取消/超时解冻回退 → 售后退回(按订单项行分摊原路回补)
-- 规则: 仅可用余额(balance-frozen的正值部分)可消费, 欠款不可
-- ============================================================

ALTER TABLE `trade_order`
  ADD COLUMN `account_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '佣金余额抵扣金额(用户资产消耗,非营销优惠;现金实付=pay_amount-account_amount,渠道单金额口径)' AFTER `point_used`;

ALTER TABLE `trade_order_item`
  ADD COLUMN `account_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '余额抵扣行分摊(合计=订单头account_amount,售后按行原路退回依据)' AFTER `point_amount`;

-- 账户流水业务类型扩位(余额消费全生命周期)
ALTER TABLE `account_log`
  MODIFY COLUMN `biz_type` TINYINT NOT NULL COMMENT '业务类型:1佣金入账 2提现冻结 3提现完成 4提现失败回退 5冲销扣回 6余额消费冻结 7余额消费完成 8余额消费退回';

-- 账户余额注释明确欠款不可消费
ALTER TABLE `user_account`
  MODIFY COLUMN `balance` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '可用余额(有符号,欠款为负,后续佣金入账抵扣;欠款不可用于余额消费)';
