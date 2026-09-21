# Feature Specification: 交易链路实现（购物车 → 订单 → 支付 → 售后）

**Feature Branch**: `（未创建）`

**Created**: 2026-09-21

**Status**: Draft

**Input**: User description: "交易链路实现——购物车/订单创建（四玩法）/支付（mock 渠道）/售后全链路 service 与 controller 实现与测试"

## User Scenarios & Testing *(mandatory)*

> 核心交易域：资金安全与库存联动的实现域。全部业务规则依据既有设计（ADR-0001 库存模型、
> 幂等总则、优惠三构成恒等式、状态机），规格只验收行为。支付渠道为 mock（真实渠道另立特性）。

### User Story 1 - 购物车（Priority: P1）

会员加购商品（重复加购数量累加）、修改数量、勾选结算、移除；结算页试算返回商品价态、四优惠构成明细（券/满减/积分/余额）、运费与现金应付，失效商品置灰提示。

**Why this priority**: 交易入口；试算是下单正确性的预演。

**Independent Test**: 加购 2 个 SKU → 改量 → 结算试算返回正确金额账本（AmountBook 各分项与合计勾稽）。

**Acceptance Scenarios**:

1. **Given** 购物车有启用 SKU，**When** 结算试算（不使用任何优惠），**Then** payAmount = 商品总额 + 运费，各优惠分项为 0。
2. **Given** 已加购的 SKU 被下架，**When** 结算试算，**Then** 该行进入 errors 提示且不计入金额。

### User Story 2 - 订单创建与管理（Priority: P1）🎯 MVP

会员创建订单（普通购物车/拼团/秒杀/砍价成交四种玩法上下文），支持幂等（request_token 重复提交返回原单）、券/积分/余额抵扣（同事务核销与冻结）；订单列表/详情（含状态时间线）；待付款可取消（释放库存+退优惠+解冻余额）；确认收货（→已完成，触发分销计提事件位）。

**Why this priority**: 交易核心；幂等与状态机是资金安全的基石。

**Independent Test**: 幂等序列（同 token 两次创建返回同单号）；取消后库存回补、券退回、余额解冻；秒杀下单扣活动分账库存。

**Acceptance Scenarios**:

1. **Given** 正确请求，**When** 同一 request_token 提交两次，**Then** 两次返回同一 orderNo，库存只锁一次。
2. **Given** 库存不足，**When** 下单，**Then** 拒绝返回 40001 且不产生订单/优惠核销等副作用。
3. **Given** 待付款订单，**When** 取消，**Then** 状态 90、库存释放、券退回、余额解冻，全部原子生效。
4. **Given** 秒杀活动剩余 1 件，**When** 两人并发下单各买 1 件，**Then** 恰好一人成功（40003 一人拒绝）。

---

### User Story 3 - 支付（mock 渠道） (Priority: P1)

会员对待付款订单发起支付（创建支付单，返回 mock 唤起参数）；模拟渠道回调驱动支付成功（库存核销、订单→待发货、通知事件位）；支持关单（订单取消/超时关闭支付单）；支付状态查询。

**Why this priority**: 资金入口；回调幂等四层防线（条件更新/唯一约束/金额校验/留档）在本故事落地。

**Independent Test**: 发起支付 → mock 回调 → 订单待发货、库存核销、余额消费完成、支付单成功（金额与订单应付+余额抵扣一致）。

**Acceptance Scenarios**:

1. **Given** 待支付订单（含余额抵扣），**When** 支付成功回调，**Then** 现金入账 = payAmount − accountAmount，余额核销完成。
2. **Given** 同一回调重放两次，**When** 处理，**Then** 第二次按重复通知直接应答成功，订单状态不再变化（幂等）。
3. **Given** 回调金额与应付不符，**When** 处理，**Then** 拒绝并留档告警。

---

### User Story 4 - 售后（Priority: P2）

已完成/待收货订单的订单项可发起售后（仅退款/退货退款，数量≤行数量）；审核通过→退款（mock 渠道，幂等 out_refund_no）→完成（回补库存、退回优惠分摊、余额退回、佣金冲销事件位、订单 refund_status 更新）；支持撤销与填写寄回单号。

**Why this priority**: 交易闭环收口；退款联动（库存/优惠/余额/佣金）是资金一致性的最后一环。

**Independent Test**: 完成订单发起退货退款 → 审核 → 确认收货 → 退款完成：库存回补、券/积分/余额按分摊回退、订单 refund_status 更新。

**Acceptance Scenarios**:

1. **Given** 已完成订单行（买 3 退 1），**When** 申请退货退款 1 件并走完流程，**Then** 库存回补 1、退款金额按行分摊、订单 refund_status 更新。
2. **Given** 待审核售后单，**When** 买家撤销，**Then** 状态 91 可再次申请。

---

### Edge Cases

- 下单幂等：同 token 不同参数的第二次提交 → 幂等冲突错误（防参数篡改）。
- 优惠券在下单与支付之间过期：已核销不回退（下单成功即占用）。
- 积分抵扣订单部分退款：积分按分摊回退（可能致负余额）。
- 拼团订单取消：释放库存 + 团成员名额释放（成员行物理删）。
- 砍价成交订单取消：砍价单状态回"到底价可下单"（可再次购买）。
- 秒杀订单取消：回补活动已售（不回普通库存）。
- 回调乱序：支付成功回调先于发起支付的响应到达（客户端视角）——状态机条件更新保证幂等。
- 售后审核与买家撤销并发：条件更新，一方成功一方失败。

## Requirements *(mandatory)*

### Functional Requirements

**购物车**

- **FR-001**: 加购 MUST 校验 SKU 可售（下架/禁售拒绝），重复加购数量累加；单行数量上限 99。
- **FR-002**: 结算试算 MUST 返回 AmountBook 全量分项（商品总额/券/满减/积分/余额/运费/现金应付）且分项与合计勾稽（恒等式）。
- **FR-003**: 试算 MUST 校验限售（收货省 ∈ 禁售列表 → 该行失效提示）。

**订单创建**

- **FR-004**: 下单 MUST 幂等：request_token 唯一约束兜底，重复提交返回原订单；同 token 不同参数 MUST 返回幂等冲突。
- **FR-005**: 下单 MUST 按序执行：限售校验 → 库存锁定（普通/拼团走 inventory；秒杀走 flash_sale_item 分账）→ 优惠计算（满减自动命中最优档 → 券 → 积分）→ 余额冻结 → 快照落库 → 状态流水。
- **FR-006**: 砍价成交 MUST 校验砍价单状态=到底价待下单，成交后砍价单置已下单（防重复成交）。
- **FR-007**: 拼团参团 MUST 校验团状态=拼团中、未超时、名额未满；人齐时置已成团（应用层触发）。

**支付**

- **FR-008**: 发起支付 MUST 创建支付单（金额=现金应付=payAmount−accountAmount），同单旧待支付支付单自动关闭。
- **FR-009**: 支付成功回调 MUST：验签(mock 为占位) → 幂等抢占（status=10 条件更新）→ 渠道单号唯一 → 金额校验 → 原文留档 → 同事务（订单 20 + 库存核销 + 余额消费完成 + 券已用确认）。
- **FR-010**: 支付失败/关单 MUST 回退余额冻结、释放库存、退回券。

**售后**

- **FR-011**: 售后申请 MUST 校验订单已完成（V1 范围）、行数量足够、金额按行分摊计算；仅退款直接进入待退款，退货退款进入待寄回。
- **FR-012**: 退款完成 MUST 同事务：库存回补（按数量）+ 优惠/余额按行分摊回退 + 订单 refund_status 更新 + 佣金冲销事件位。

**订单管理（admin）**

- **FR-013**: 管理员发货（状态 20→30，记录物流）、取消、卖家备注 MUST 走状态机条件更新并写状态流水与操作留痕。

### Key Entities *(include if feature involves data)*

既有表承载（零新表）：trade_order / trade_order_item / trade_order_log / cart_item / pay_order / pay_callback_log / after_sale_order / inventory / inventory_log / flash_sale_item / bargain_record / group_buy_team(_member) / user_coupon / point_account(_log) / user_account / account_log / promotion_*（V6~V25/V28/V29）。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 恒等式全场景成立：`promotion_amount ≡ coupon + full_reduction + point`（头/项两层）；`pay_order.amount ≡ pay_amount − account_amount`；`pay_amount ≡ total − promotion + freight`。
- **SC-002**: 幂等：同 token 重复下单返回同单号且库存/优惠零重复扣减。
- **SC-003**: 状态机防御：全部非法迁移（已完成再取消、待发货直接确认等）返回 40006，零例外。
- **SC-004**: 库存一致性：下单锁定/支付核销/取消释放/售后回补四操作后，`inventory` 与 `inventory_log` 快照对账一致。
- **SC-005**: 回调幂等：同一支付回调重放，第二次零状态变化；金额不符被拒并留档。
- **SC-006**: 余额消费资金链完整：冻结→核销→（取消）解冻→（售后）退回，`account_log` 双向可追溯且余额勾稽。
- **SC-007**: 全部 service 方法红绿 TDD 覆盖（购物车/下单/支付/售后四大块），`make build && make test` 全绿。

## Assumptions

- 支付渠道为 mock（内部直接驱动回调），真实微信支付另立特性；回调验签为占位接口。
- 短信/通知在订单状态迁移处只创建通知任务位（通知投递属消息域特性）。
- 评价、分销计提的真实业务在各自特性实现，本特性只预留事件调用位（接口注释标注）。
- 余额提现、佣金结算属分销域（已设计不在本特性）。
- 运费计算复用 005 域的 FreightTemplate 口径（V1 默认模板可空=包邮）。
- 单测数据策略延续 005：自建+前置清理，确定性配置注入。
