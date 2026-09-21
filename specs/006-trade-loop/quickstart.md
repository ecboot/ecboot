# Quickstart: 交易链路验证指南

前提：compose 全栈 + `make migrate-up` + 003 认证可用（会员 token）+ 测试商品数据（TDD 自建或 quickstart 手工建）。

## 场景一：下单幂等 + 金额勾稽（SC-001/002）

```bash
BASE=http://127.0.0.1:8080 TOKEN=<会员token>
# 试算
curl -s "$BASE/shop/cart/checkout?addressId=1" -H "Authorization: Bearer $TOKEN"
# 下单（同 token 提交两次）
curl -s -X POST "$BASE/shop/orders" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"requestToken":"tk-001","addressId":1,"cartItemIds":[1]}'
# 预期: 两次返回同一 orderNo; AmountBook 勾稽:
#   payAmount ≡ totalAmount − couponAmount − fullReductionAmount − pointAmount + freightAmount
```

## 场景二：库存一致性（SC-004）

下单前后 `SELECT total,locked FROM inventory WHERE sku_id=?`：locked +n；取消后 locked 回落；支付后 total −n。inventory_log 每步一条（change_type 1/2/3）。

## 场景三：支付幂等 + 余额（SC-005/SC-006）

```bash
curl -s -X POST "$BASE/shop/pay" -H "Authorization: Bearer $TOKEN" -d '{"orderNo":"..."}'
# mock 回调（测试驱动或 mock 端点）两次:
#   第一次: 支付单 20, 订单 20, 余额 frozen 扣减, account_log bizType=7
#   第二次: 直接幂等应答, 零状态变化
# 余额恒等式: payOrder.amount ≡ trade_order.pay_amount − account_amount
```

## 场景四：售后回退（FR-011/012）

申请退货退款（数量 1）→ 审核 → 确认收货 → 退款完成（mock 回调）：库存回补 1、券/积分/余额按行分摊回退、订单 refund_status 更新、account_log bizType=8。

## 宪法 IV

`cd apps/server && make build && make test`（交易域集成测试全绿）。
