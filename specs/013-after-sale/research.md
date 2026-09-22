# Phase 0 研究与决策：售后域（013-after-sale）

**Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

本批性质：**接口与状态机已定义、impl 待补 + 11 个 controller 桩连线**。因此研究集中在"既有契约里没写清、或写错"的几处，逐一裁定并给出依据。

## D1 退款发起时机与状态推进

- **Decision**: 由"进入待退款"的同一个动作**立即尝试发起渠道退款**——`Approve`（仅退款）与 `ConfirmReceipt`（退货退款）在把状态推到 30 后立刻调用渠道；成功则推进到 40（退款中），失败则**停留在 30** 并记录失败原因，等待后台 `RetryRefund`。
- **Rationale**: 11 个方法里没有独立的"发起退款"语义（`Approve`/`ConfirmReceipt`/`RetryRefund` 三个入口已覆盖全部发起时机）；30 与 40 的区分正是"未发起/在途"，若审核通过后不自动发起，30 会长期悬挂且没有任何端点能推进它。
- **Alternatives considered**: ① 引入定时任务扫描 30 并发起（新增调度面，且本批无调度需求）；② 增设"发起退款"端点（改契约，违背"接口已有"前提）。

## D2 退款渠道的幂等键与失败语义

- **Decision**: 沿用 `out_refund_no = after_sale_no`（既有契约注释 + 迁移 000008 列注释）；渠道失败**不停留在"已完成"**，而是 30 + 失败原因；`RetryRefund` 允许 30 → 重试，已完成（50）不可再退。
- **Rationale**: `paychannel.Channel.Refund(outRefundNo, amountFen)` 的注释即"幂等: outRefundNo"，与售后单号天然对齐；回调侧 `HandleRefundNotify`（批次 06 已修）只认 `after_sale_no` 且只在 40 时推进到 50，因此"我方重试 + 回调确认"不会重复出款。
- **Alternatives considered**: 用 `after_sale_no + 序号` 作多退单号——本批每单一退，无需。

## D3 审核/退款操作人的审计载体（**需补迁移**）

- **Decision**: 新增迁移 **000038_after_sale_operator**，给 `after_sale_order` 加两列：`operator_id VARCHAR(64) NOT NULL DEFAULT ''`（最后操作人 `admin:{id}`，覆盖审核/确认收货/退款重试）与 `fail_reason VARCHAR(255) NOT NULL DEFAULT ''`（渠道退款失败原因）。`gf gen dao` 重新生成对应 entity/do。
- **Rationale**: `IAfterSaleLogic` 的 `Approve/Reject/ConfirmReceipt/RetryRefund` **都接收 operator 参数，而表里没有任何列能承载它**——不补就是"审计信息被静默丢弃"，正是批次 06 评审在 `AdminCancel` 抓到的同型缺陷（Important）。`admin_operation_log` 是请求级审计表但**至今无人写入**（其列表端点属批次 12 的风控审计范围），不能作为本批的载体。
- **Alternatives considered**: ① 不补、记已知缺口（重蹈覆辙，否决）；② 复用 `reject_reason`/`description` 等业务列（语义污染，否决）；③ 新建 `after_sale_log` 表（本批只需"最后操作人"，一表过重，且状态流水在 50 前只有 4 跳）。
- **代价与记账**: 本批因此**不是零迁移**（spec 的 Assumptions 已在 Phase 1 复查时同步修正）。已有先例：批次 01 加 000035、批次 05 加 000036，均按协议记账。

## D4 撤销的状态边界（用户裁定）

- **Decision**: 可撤 = **10 待审核 / 20 待买家寄回 / 30 待退款**；**40 退款中及以上不可撤**（含各终态）。
- **Rationale**: 在途退款不可中断。若允许撤销 40，渠道退款成功后回调会因状态不符（≠40）而不匹配，形成"钱已出、单已撤销、对账捞不到"的黑洞——与批次 06 的 C3a 静默资损面同型。用户于 2026-09-22 明确选择本方案（选项 B）。
- **Alternatives considered**: ① 仅 10 可撤（审核通过后买家无法自主取，客服工单上升）；② 未终态全可撤（沿用原注释，但需额外建立"撤销后在途退款对账标记"，成本与风险更高）。
- **口径修正**: `IAfterSaleLogic.Cancel` 的原注释"未终态→91撤销"据此收窄，实现时同步注释，避免下轮照旧注释写回（该教训在批次 06 已发生过一次：图纸未改导致 30/90 口径反复）。

## D5 "可退数量"的累计口径与拒绝/撤销的释放

- **Decision**: 可退数量 = 订单项购买数量 − **该行上"未被释放"的售后数量之和**；"未被释放" = 状态不属于 {已拒绝 90, 已撤销 91}。
- **Rationale**: 防超额必须是累计口径（本批 spec 的 FR-002 与 SC-003）。而已拒绝/已撤销若不释放，一笔误操作就会把该行永久锁死——与"拒绝即释放"的用户预期相悖。
- **Alternatives considered**: 拒绝即保留额度（更保守但体验差，且与"可再次申请"的常规电商行为不符）。

## D6 退款金额口径与 0 元边界

- **Decision**: 退款金额由系统计算：`行实付金额 × quantity ÷ 行数量`，向下取整到分，**最后一笔退完剩余行数量时取剩余全额**（保证多次部分退之和恰好等于行实付，无尾差溢出）；若计算结果为 **0**（全额优惠的行）→ **不调用渠道**，直接以"已完成"收口并留痕（`fail_reason` 为空、`refund_no` 为空）。
- **Rationale**: `refund_amount ≤ 订单项 pay_amount`（000008 列注释）必须成立；批次 06 已把行实付（含优惠分摊）写进 `trade_order_item.pay_amount`，本批直接消费该口径，无需重新分摊。0 元退款调用渠道会被渠道拒绝或产生无意义流水。
- **Alternatives considered**: 允许买家指定金额（api 契约无此字段，否决）；按比例四舍五入（多次部分退可能超过行实付，风险更大，否决）。

## D7 完成后的库存回补与订单退款状态

- **Decision**: 仅**退货退款**在完成时按 `quantity` 回补 `inventory.total`（`locked` 不动）；随后重算订单 `refund_status`：订单全部行的"已售后退款数量"之和 == 全部行购买数量之和 → 2 全额退款；>0 且不等 → 1 部分退款；否则 0。（"已退款数量"只统计状态=50 已完成的行。）
- **Rationale**: 仅退款货物未退回，回补即虚增可售（超卖风险）；`refund_status` 的注释为"0无售后 1部分退款 2全额退款"，按"完成态"统计才与买家感知一致。
- **Alternatives considered**: 状态 40 即计入退款状态（在途未定，易与失败重试打架，否决）。

## D8 佣金冲销事件出口

- **Decision**: 售后完成时调用既有 `shop.NotifyEnq`（`INotifyEnqueue`）之外的**独立出口**：本批只落一条"佣金冲销待处理"的意图记录（复用既有事件位风格），结算与账务处理属批次 11（分销与资金）。若批次 11 的结算按"订单完成态"计提，则冲销必须是**显式事件**才能被消费。
- **Rationale**: 避免"退了款但佣金仍计提"。批次 06 已建立"事件位留出口、结算后置"的先例（`Confirm` 只落状态与流水，计提由分销结算推进）。
- **Alternatives considered**: 本批直接冲销佣金账户（依赖账户域，批次 11 才存在，越界）。

## D9 复用既有幂等与并发防线（不重复造）

- **Decision**: 全部状态迁移一律"条件更新 + 判 RowsAffected"（复用批次 06 修正后的写法）；退款回调直接复用已修的 `HandleRefundNotify`（`after_sale_no` 匹配 + 40→50 + 失败留档 + affected=0 不记已处理）。
- **Rationale**: 批次 06 的评审与修复已把这两处打磨到"条件更新必判行数"的既有惯例；新代码若再出现"只看 error"的写法就是已知缺陷复发（I1/I11 同型）。
- **Alternatives considered**: 无（这是既有工程惯例，不存在可选方案）。

## 未知项收敛

| Technical Context 中的未知 | 结论 |
|---|---|
| 是否新增迁移 | **是**，唯一迁移 000038（D3），两列追加 |
| 退款渠道调用方式 | 复用 `paychannel.Channel.Refund`（D2） |
| 操作人承载 | 新增 `operator_id` 列（D3） |
| 撤销边界 | 用户裁定（D4） |
| 0 元退款 | 直收口不调渠道（D6） |

**结论**：无遗留 NEEDS CLARIFICATION，可进入 Phase 1。
