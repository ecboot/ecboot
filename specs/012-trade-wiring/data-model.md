# Data Model: 012-trade-wiring

> 零表结构新增。本文档摘录本批消费的字段与语义（权威以迁移与 gf 生成物为准）。

## 一、订单（trade_order / trade_order_item / trade_order_log）

**状态机**：10 待付款 → 20 待发货 → 30 待收货 → 40 已完成；10 → 90 已取消。

本批新增操作：
| 操作 | 迁移 | 前置条件 | 副作用 |
|---|---|---|---|
| Confirm（确认收货） | 30 → 40 | 仅状态 30（否则 40006） | 订单日志 + 佣金计提事件位 |
| Deliver（发货） | 20 → 30 | 仅状态 20；物流公司存在且启用 | 落 deliver_company/deliver_no + 日志 + 通知事件位 |
| AdminCancel | → 90 | 非终态 | 同 C 端取消（释放库存/券/余额/积分） |
| SellerRemark | 不变 | 任意 | 写 seller_remark（C 端 Detail 查询不含该列） |
| CancelTimeout / AutoConfirm | 定时 | 超时/发货后 N 天（配置化） | 复用取消/确认逻辑 |

## 二、支付（pay_order / pay_callback_log）

### pay_order
| 字段 | 语义 | 本批规则 |
|---|---|---|
| pay_no | 支付号 | 幂等键 |
| order_no / order_id / user_id | 关联 | — |
| pay_channel | 渠道 | mock（006 的 paychannel） |
| amount | 应付金额（decimal） | **回调金额校验基准** |
| status | 10 待支付 / 20 成功 / 30 支付失败 / 90 已关闭 | 条件更新推进（口径以 000007 列注释为准；修复轮勘误: 原图纸误写"30 关闭"） |
| channel_trade_no / success_time / closed_time / fail_reason | 渠道与结果 | 回调填 |

### pay_callback_log（只追加，对账依据）
渠道 + 原文 + 处理结果；**成败均留档**（重放识别与审计）。

### 回调幂等四层防线（FR-006/007）
```
① 条件更新: UPDATE pay_order SET status=20 WHERE pay_no=? AND status=10
   affected=0 → 已处理/不存在（直接应答成功, 幂等）
② 唯一约束: pay_callback_log 落档（同渠道+流水重放可识别）
③ 金额校验: 回调金额 vs pay_order.amount（不符 → 留档 + 拒绝）
④ 原文留档: 无论成败均落原文
同事务推进: 支付单 20 → 订单 20（条件 status=10）→ 库存核销 → 余额消费完成
脏态保护: 订单已取消（90）→ 不推进订单（留档）
```

## 三、优惠券（coupon / user_coupon）

### coupon（模板）
threshold_amount（门槛）、discount_amount（抵扣值）、discount_rate（折扣率）、total_count/received_count（防超发）、
per_limit（每人限领）、valid_type/valid_start_at/valid_end_at/valid_days（有效期三形态）、status/deleted。

### user_coupon（会员券）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| status | 1 未使用 / 2 已使用 / 3 已过期 / 4 已退回 | **过期惰性判定**（查询时比较 expire_time, 不改表） |
| expire_time | 到期时间 | 领券时按模板有效期计算落库 |
| order_no / used_time | 核销绑定 | Consume 时写；ReturnBack 时清除绑定（有效期不变） |

### 领取（FR-013）同事务
```
① 限领校验: COUNT(user_coupon WHERE user_id=? AND coupon_id=?) < per_limit
② 防超发: UPDATE coupon SET received_count=received_count+1
          WHERE id=? AND received_count < total_count   （affected=0 → 30001/50001 已领完）
③ 计算 expire_time（valid_type 三形态）→ INSERT user_coupon(status=1)
```

## 四、积分抵扣修复（FR-001）

`calcPointDeductFen`：查询移除 `AND deleted=0`（该表无此列）；余额<=0 或未勾选 → 0；无账户行 → 0（不报错）。
算法不变：`min(balance, maxFen)`（1 积分=1 分）。

## 五、DTO 与迁移

- **DTO 零新增**（既有 100 类型覆盖）。
- **零迁移**（全部表既有）。
