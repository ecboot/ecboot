# Phase 1 数据模型：售后域（013-after-sale）

**依据**: 迁移 `000008_after_sale_domain`（建表）+ `000023_business_review_fixes`（补 quantity/currency）+ **本批新增 `000038_after_sale_operator`**（两列追加）。

## 一、after_sale_order（既有表 + 本批两列）

| 列 | 类型 | 说明（契约） |
|---|---|---|
| id / after_sale_no | BIGINT / VARCHAR(32) | 主键 / 售后单号（唯一；**兼作渠道退款 out_refund_no**） |
| order_id / order_no | BIGINT / VARCHAR(32) | 所属订单（冗余单号免联查） |
| order_item_id | BIGINT | 售后粒度 = **订单项**（一单一项） |
| user_id | BIGINT | 申请人 |
| type | TINYINT | 1 仅退款 / 2 退货退款 |
| status | TINYINT | 见 §二 状态机 |
| currency | CHAR(3) | 随订单（默认 CNY） |
| quantity | INT UNSIGNED | 退货/退款**数量**（≤ 订单项数量） |
| reason / description | VARCHAR | 售后原因 / 问题描述 |
| voucher_images | JSON | 凭证图片数组 |
| refund_amount | DECIMAL(10,2) | 退款金额（≤ 订单项 pay_amount，系统计算） |
| return_logistics_no | VARCHAR(64) | 买家寄回单号（type=2） |
| refund_no | VARCHAR(64) | 渠道退款单号（回填；mock 渠道下可为空） |
| reject_reason | VARCHAR(255) | 拒绝原因（对买家可见） |
| audit_time | DATETIME | 审核时间 |
| refund_time | DATETIME | **退款完成**时间 |
| **operator_id** | VARCHAR(64) NOT NULL DEFAULT '' | **新增**：最后操作人 `admin:{id}`（审核 / 确认收货 / 退款重试） |
| **fail_reason** | VARCHAR(255) NOT NULL DEFAULT '' | **新增**：渠道退款失败原因（重试前保留，成功后清空） |
| created_at / updated_at | DATETIME | 既有 |

**为什么加这两列**：`IAfterSaleLogic` 的四个后台方法都接收 `operator`，而原表无列可承载 → 审计信息会被静默丢弃（批次 06 评审同型缺陷）；`fail_reason` 让"待退款 + 失败原因 + 可重试"这条路径有据可查（FR-012）。

## 二、状态机（唯一事实源）

```
                    ┌─ 仅退款 ──────────────┐
10 待审核 ─(同意)───┤                       ├─→ 30 待退款 ─(发起退款成功)→ 40 退款中 ─(渠道回调)→ 50 已完成
                    └─ 退货退款 → 20 待买家寄回 ─(确认收货)─┘         │
                                                                      └─(发起失败)→ 30（记 fail_reason，可 RetryRefund）
10 ─(拒绝，必填原因)→ 90 已拒绝
可撤集合 {10, 20, 30} ─(买家撤销)→ 91 已撤销
```

| 迁移 | 允许的前态 | 触发 | 副作用 |
|---|---|---|---|
| → 10 | （新建） | 会员申请 | 落成因快照 + 计算 refund_amount |
| 10 → 20 | 10 | 后台同意（type=2） | audit_time、operator_id |
| 10 → 30 → 40 | 10（type=1） | 后台同意后**立即发起退款** | audit_time、operator_id；成功→40，失败→留在 30 + fail_reason |
| 20 → 30 → 40 | 20（已填寄回单号） | 后台确认收货后**立即发起退款** | 同上 |
| 30 → 40 | 30 | 后台重试退款 | 清 fail_reason；成功→40 |
| 40 → 50 | 40 | **渠道退款回调**（复用 `HandleRefundNotify`，键=`after_sale_no`） | refund_time |
| 10 → 90 | 10 | 后台拒绝（原因必填） | reject_reason、audit_time |
| {10,20,30} → 91 | 10/20/30 | 买家撤销 | （释放可退数量） |

**铁律**：所有迁移一律**条件更新 + 判 RowsAffected**（affected=0 → 40006 状态不允许），与批次 06 修正后的写法一致；40 及以上不可撤销（在途资金不可中断，见 research D4）。

## 三、退款金额计算（D6）

```
行数量 = N, 行实付 = P（trade_order_item.pay_amount，含行优惠分摊）
本次申请数量 = q
本行已完成的退款数量 = q_done（仅统计状态=50 的单）
若 q_done + q == N  →  refund_amount = P - Σ(已完成的 refund_amount)   // 末笔取剩余全额，消尾差
否则                →  refund_amount = floor(P * q / N)                // 向下取整到分
约束: refund_amount ≤ P（超额即拒绝 40008）
边界: refund_amount == 0 → 不调用渠道，直接以已完成收口（记 refund_time，refund_no 留空）
```

## 四、辅助实体（只读消费）

- **trade_order_item**：提供 `quantity` / `pay_amount`（行实付）/ `sku_id` —— 申请校验与金额计算依据；**不改写入**。
- **trade_order**：完成后按"全部行的已完成退款数量"更新 `refund_status`（0 无售后 / 1 部分退款 / 2 全额退款）。
- **inventory**：**仅退货退款**完成时 `total += quantity`（`locked` 不动；仅退款不回补）。
- **事件出口**：完成时投递"佣金冲销"意图（结算属批次 11）。

## 五、校验规则汇总（可测）

| 规则 | 出处 | 失败码 |
|---|---|---|
| 订单项归属本人 | FR-001/002 | 40009（按不存在处理，不泄露他人资源） |
| 订单项所属订单可售后（已完成 40；已发货未完成是否可退见 spec Assumptions） | FR-002 | 40008 |
| 1 ≤ quantity ≤ 行数量 | FR-002 | 40008 |
| 累计可退数量足够（不释放 90/91） | FR-002/SC-003 | 40008 |
| refund_amount ≤ 行实付 | FR-003 | 40008 |
| 状态机合法性（每跳前态匹配） | FR-007~013 | 40006 |
| 后台拒绝必填原因 | FR-008 | 10001（参数校验） |
| 确认收货须已填寄回单号 | FR-010 | 40008 |
| 状态迁移并发不重入 | FR-017 | 40006（affected=0） |
