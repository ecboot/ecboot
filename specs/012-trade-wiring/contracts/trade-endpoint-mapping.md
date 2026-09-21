# Contract: 012-trade-wiring 端点 × 方法（22 端点）

> 调用形态：shop 域 struct（`shop.NewOrderLogic()` / `shop.NewPayLogic()`）、user 域包级函数。
> 越权：C 端订单/券方法以会话 userId 收口（他人订单按"不存在"）；后台方法无 userId（运营视角，权限点控制）。

## shop 渠道（14）

### 购物车（5，连线既有）
| 端点 | 方法+路径 | service |
|---|---|---|
| CartDetail | GET /shop/cart | `CartLogicImpl.Detail` |
| CartAddItem | POST /shop/cart/items | `CartLogicImpl.AddItem` |
| CartUpdateItem | PUT /shop/cart/items/{id} | `CartLogicImpl.UpdateItem` |
| CartRemoveItem | DELETE /shop/cart/items/{id} | `CartLogicImpl.RemoveItem` |
| CartCheckout | POST /shop/cart/checkout | `CartLogicImpl.Checkout`（含 FR-001 修复后积分抵扣） |

### 订单 C 端（5）
| 端点 | 方法+路径 | service |
|---|---|---|
| OrderCreate | POST /shop/orders | `OrderLogicImpl.Create`（幂等 token） |
| OrderList | GET /shop/orders | `OrderLogicImpl.List` |
| OrderDetail | GET /shop/orders/{orderNo} | `OrderLogicImpl.Detail` |
| OrderCancel | POST /shop/orders/{orderNo}/cancel | `OrderLogicImpl.Cancel` |
| OrderConfirm | POST /shop/orders/{orderNo}/confirm | `OrderLogicImpl.Confirm`（**新写**） |

### 支付（4，**新写**）
| 端点 | 方法+路径 | service |
|---|---|---|
| PayCreate | POST /shop/pay/create | `PayLogicImpl.Create` |
| PayStatus | GET /shop/pay/status | `PayLogicImpl.Status` |
| PayNotify | POST /shop/pay/notify | `PayLogicImpl.HandlePayNotify`（公开：渠道签名语义） |
| RefundNotify | POST /shop/refund/notify | `PayLogicImpl.HandleRefundNotify`（公开） |

## admin 渠道（5，**新写**）

| 端点 | 方法+路径 | service | 权限点 |
|---|---|---|---|
| AdminOrderList | GET /admin/orders | `OrderLogicImpl.AdminList` | 登录 |
| AdminOrderDetail | GET /admin/orders/{orderNo} | `OrderLogicImpl.AdminDetail` | 登录 |
| AdminOrderDeliver | POST /admin/orders/{orderNo}/deliver | `OrderLogicImpl.Deliver` | order:deliver |
| AdminOrderCancel | POST /admin/orders/{orderNo}/cancel | `OrderLogicImpl.AdminCancel` | order:cancel |
| AdminOrderRemark | POST /admin/orders/{orderNo}/remark | `OrderLogicImpl.SellerRemark` | order:update |

## user 渠道（3，**新写**）

| 端点 | 方法+路径 | service |
|---|---|---|
| CouponAvailableList | GET /user/coupons/available | `AvailableTemplates`（包级） |
| CouponReceive | POST /user/coupons/{id}/receive | `Receive` |
| MyCouponList | GET /user/coupons | `Mine`（status 筛选 + 惰性过期） |

## 内部方法（实现, 不建端点）
- `OrderLogicImpl.CancelTimeout` / `AutoConfirm`（定时任务）
- `UsableForOrder` / `Consume` / `ReturnBack`（结算匹配/核销/退回——订单与结算编排调用）
- `PayLogicImpl.CloseExpired`（关单扫描）
- 前置修复：`calcPointDeductFen`（集成在 Checkout，无端点变化）

## 错误码（复用既有）
| 码 | 触发 |
|---|---|
| 40001 | 库存不足（既有） |
| 40005 | 订单不存在（含他人订单——不可见语义） |
| 40006 | 状态不允许该操作（确认收货/发货/取消的前置态不符） |
| 40007 | 支付单创建失败（应付<=0 等） |
| 50001 | 券已领完/超限（既有） |
| 10001/10006 | 参数非法 / 资源不存在（含物流公司停用→拒发） |
