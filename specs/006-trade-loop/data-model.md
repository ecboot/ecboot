# Data Model: 交易链路（Phase 1）

**零新表**。数据形态 = 既有交易表的写路径编排 + 事务边界 + 跨域接口锁定。

## 一、下单事务编排（order_impl.Create，唯一事务入口）

```text
Begin
 1. 幂等抢占: INSERT trade_order(request_token 唯一) —— 1062 → 反查校验指纹 → 返回原单
 2. 限售校验: 收货省 ∉ 各行 SPU 禁售列表
 3. 库存锁定:
    普通/拼团 → inventory.Lock(tx, skuId, qty, orderNo)
    秒杀      → flash_sale_item 条件更新(sold_count+n WHERE stock-sold>=n)
 4. 优惠计算: 满减(范围内商品金额命中档) → 券核销(状态/门槛/过期校验) → 积分抵扣 → 余额冻结
 5. 分摊: 四构成按行金额比例, 尾差记末行(各构成独立分摊)
 6. 快照落库: trade_order(金额账本+收货+归因+玩法关联) + trade_order_item(商品快照+行分摊)
 7. 余额冻结: user_account balance→frozen + account_log(bizType=6)
 8. 玩法联动: 拼团成员行插入(人齐置成团) / 秒杀已售已增 / 砍价单置已下单
 9. 状态流水: from=NULL,to=10
Commit（任一步失败全量回滚; Redis 幂等 token 由调用方清理）
```

## 二、支付成功事务（pay_impl.HandlePayNotify）

```text
Begin
 1. 支付单抢占: UPDATE pay_order SET status=20 WHERE pay_no=? AND status=10 —— affected=0 → 幂等应答
 2. 渠道单号唯一校验; 金额校验: 回调金额 ≡ pay_order.amount（现金应付）
 3. 订单: status 10→20（条件更新）; pay_time
 4. 库存核销: locked 行 inventory.Deduct(tx, ...)
 5. 余额消费完成: frozen 扣减 + account_log(bizType=7)
 6. 券已使用确认（下单时已核销, 此处无操作——留注释）
 7. 通知事件位: 订单已支付（notify 域）
 8. pay_callback_log 原文留档
Commit
```

## 三、跨域接口锁定（依赖倒置, user 域注册实现）

| 接口（shop 域定义） | 方法 | 实现域 |
|---|---|---|
| ICouponTrade | ConsumeForOrder / ReturnForOrder（tx 感知） | user |
| IPointTrade | ConsumeForOrder / RefundForOrder（tx 感知） | user |
| IAccountTrade | FreezeForOrder / SettleFrozen / Unfreeze / RefundFrozen（tx 感知） | user |
| INotifyEnqueue | OrderPaid / OrderShipped / ...（事件位） | user |

注册方式：`service/shop` 定义接口变量 + `RegisterXxx(impl)`（user 域 bootstrap 装配时注入）——依赖倒置，shop 不 import user。

## 四、状态机迁移实现（全部条件 UPDATE）

| 迁移 | SQL 语义 | 副作用（同事务） |
|---|---|---|
| 10→20 | 支付回调驱动 | 库存核销+余额完成+事件 |
| 10→90 | 条件 status=10 | 释放库存+退券(退回)+积分回退+余额解冻+砍价单回退+拼团名额释放 |
| 20→30 | 发货 | 物流信息+流水 |
| 30→40 | 确认收货/7天自动 | finish_time+分销计提事件位 |

非法迁移：affected=0 → 40006。

## 五、金额口径（int64 分）

- `money.FromYuanString(s) (int64, error)` / `money.ToYuanString(f int64) string`（`internal/library/money` 新建，分转元两位小数）。
- 恒等式断言全部在分域：
  - `promotion ≡ coupon + fullReduction + point`（头/项）
  - `pay ≡ total − promotion + freight`
  - `payOrder.amount ≡ pay − account`
