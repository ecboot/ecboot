# Research: 012-trade-wiring

> Phase 0 产出。基于代码勘察（commit 3d39a2f 时点）。

## D1 前置修复：积分抵扣（批次 05 评审发现的跨批次 bug）

**结论**：`promotion_calc.go` 的 `calcPointDeductFen` 查询
`SELECT balance FROM point_account WHERE user_id=? AND deleted=0` → **去掉 `AND deleted=0`**。

**理由**：`point_account` 表无 `deleted` 列（V19 建表 + V26 仅加 `last_earned_at`，全局 grep 确认）→
SQL 必报错 → `err != nil` 静默返 0 → **积分抵扣恒为 0**。修复后按既有算法
（1 积分=1 分；min(余额, 可抵上限)）生效；`balance<=0` 仍返 0（负余额不抵扣 ✓）。

**验证**：新增行为级断言（有余额勾选→>0；无余额/未勾选→0；无账户行→0 不报错）。

## D2 支付实现：幂等四层防线（照 006 契约）

**结论**：`HandlePayNotify` 按接口注释实现四层防线：
1. **条件更新**：`UPDATE pay_order SET status=20 WHERE pay_no=? AND status=10`（affected=0 → 已处理或不存在）
2. **唯一约束**：`pay_callback_log` 留档（渠道+流水号唯一）——重放识别
3. **金额校验**：回调金额 vs `pay_order.amount`（不符 → 留档 + 拒绝）
4. **原文留档**：无论成败均落原文（对账依据）

同事务推进：支付单 20 → 订单 20（条件：订单状态 10 → 20）→ 库存核销（Deduct）→ 余额消费完成
（account 域 ports，006 已建立）。

**条件更新用 `RowsAffected` 判定**（gf 未开 clientFoundRows，返回变更行数）✓（010 已验证同机制）。

## D3 后台订单管理：追加到既有 OrderLogicImpl

**结论**：`Confirm`/`AdminList`/`AdminDetail`/`Deliver`/`AdminCancel`/`SellerRemark`/`CancelTimeout`/`AutoConfirm`
追加为 `OrderLogicImpl` 的方法（不新建结构——宪法 V）。复用既有私有 helper（如取消的释放逻辑）。

**发货**：状态机 20→30（条件更新）+ 物流信息落列 + **物流公司校验**（`logistics_company` 存在且 status=1，
批次 03 字典）+ 订单日志。
**取消**：与 C 端 Cancel 同语义（复用其释放逻辑：库存/券/余额/积分）。
**备注**：写 `trade_order.seller_remark`（V28 列），C 端 Detail 的 Fields 不含该列 ✓（买家不可见由查询口径保证）。

## D4 优惠券：惰性过期 + 防超发 + 限领

**结论**：
- `Mine`：四态筛选；**已过期惰性判定**——查询时 `WHERE status=1 AND expire_at < NOW()` 视为"已过期"展示
  （不改表、不依赖定时任务）
- `Receive`：同事务「条件更新模板已领数 +1（防超发: `WHERE received < total` 若有该列；否则按既有列）」
  + 限领校验（个人已领数 < 每人限领）+ 插 user_coupon
- `UsableForOrder`：门槛 ≤ 商品金额 + 未使用 + 未过期 → 按抵扣降序
- `Consume`：置已使用 + 绑单号（下单事务内由订单域调用）
- `ReturnBack`：置未使用（有效期不变）

**具体列名以实现期读表为准**（coupon/user_coupon，000009 建表 + 后续迁移）。

## D5 连线形态

**结论**：shop 域用 struct（`NewOrderLogic()`/`NewPayLogic()`——与既有 `OrderLogicImpl` 一致）；
user 域用包级函数（`user.CouponMine(...)`——与 011 一致）。

## 勘察结论（非决策）

- **impl 覆盖**：`CartLogicImpl`（5 方法 ✓）、`OrderLogicImpl`（C 端 4/5——缺 Confirm）；`IPayLogic` 与
  `IUserCouponLogic` 无实现。
- **DTO 齐备**：`dto_shop.go` 76 + `dto_user.go` 24 类型含本批全部出入参（PayCreated/PayStatus/
  AdminOrderQuery/AdminOrderSummary/AvailableCoupon/MyCouponItem/UsableCouponItem 等）。
- **零迁移**：trade_order/pay_order/pay_callback_log/coupon/user_coupon 全既有（000006/000007/000009/000025）。
- **既有资产**：006 的 `ports.go`（跨域接口：库存/券/积分/余额）、`trade_test.go` fixture 可直接复用。
