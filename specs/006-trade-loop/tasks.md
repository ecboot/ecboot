---
description: "Task list for 交易链路实现（trade-loop）"
---

# Tasks: 交易链路实现（trade-loop）

**Input**: Design documents from `/specs/006-trade-loop/`（spec/plan/research/data-model/quickstart）

**Prerequisites**: 特性 005 已合入（IProductLogic/IInventoryLogic 及商品域可用）；003 认证可用。

**Tests**: TDD 红绿强制（宪法 §五）——交易域集成测试 `internal/service/shop/trade_test.go`：确定性配置注入、数据自建+前置清理、恒等式/幂等/状态机/库存对账断言（SC-001~006）。每故事"先红后绿"。

**Organization**: 按 spec 用户故事 US1~US4；落位 `internal/service/shop` 四实现文件 + `internal/library/{money,paychannel}`。

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 新建 `internal/library/money/money.go`：元 string ↔ 分 int64 转换（`FromYuanString`/`ToYuanString`/`Add`/`Sub`/`MulQty`），解析失败返回错误；单测覆盖精度边界（0.1+0.2 类、负数、超大额）
- [ ] T002 扩表生成：`gf gen dao -t "trade_order,trade_order_item,trade_order_log,cart_item,pay_order,pay_callback_log,after_sale_order,user_coupon,point_account,point_log,user_account,account_log,flash_sale_item,group_buy_team,group_buy_team_member,bargain_record,product_sku,product_spu,user"`（一次生成交易域全部依赖表）
- [ ] T003 [P] `internal/library/paychannel/mock.go`：PayChannel 接口（CreateParams/VerifyNotify/ParseAmount）+ Mock 实现（生成占位唤起参数、构造标准回调报文）

**Checkpoint**: 金额库/dao/mock 渠道就绪

---

## Phase 2: Foundational (Blocking Prerequisites)

- [ ] T004 跨域接口锁定：`internal/service/shop/ports.go` 定义 `ICouponTrade`（ConsumeForOrder/ReturnForOrder）、`IPointTrade`（ConsumeForOrder/RefundForOrder）、`IAccountTrade`（FreezeForOrder/SettleFrozen/Unfreeze/RefundFrozen）、`INotifyEnqueue`（OrderPaid/OrderShipped/OrderCancelled）——接口变量 + `RegisterXxx` 注入函数（依赖倒置, research D1/plan §跨域接口锁定）
- [ ] T005 交易域测试基座 `internal/service/shop/trade_test.go`：确定性配置注入（gdb/gredis SetConfig）、测试商品/地址/库存/账户数据 builder（自建+前置清理）、恒等式断言 helper（三条恒等式 + 余额式，data-model §五）

**Checkpoint**: 事务编排便可开始

---

## Phase 3: User Story 1 - 购物车（Priority: P1）

**Goal**: 加购/改量/勾选/移除/结算试算（AmountBook 全分项）

**Independent Test**: 加购 2 SKU → 改量 → 试算：无优惠时 payAmount = total + freight；失效行进 errors；四优惠构成零值勾稽

- [ ] T006 [US1] 先写失败测试 `internal/service/shop/cart_impl_test.go`：加购累加/上限 99/改量勾选移除/试算金额勾稽/失效行 errors
- [ ] T007 [US1] 实现 `internal/service/shop/cart_impl.go`：IcartLogic 全方法 + Checkout 试算（可售联查、限售校验、满减命中最优档、券匹配、积分抵扣试算、运费计费），AmountBook 分项勾稽（分域 int64）
- [ ] T008 [US1] 控制器接线：`internal/controller/shop` cart 相关桩填充（Detail/AddItem/UpdateItem/RemoveItem/Checkout）+ 转绿

**Checkpoint**: 试算与购物车契约闭环

---

## Phase 4: User Story 2 - 订单创建与管理 (Priority: P1) 🎯 核心

**Goal**: 下单九步事务编排（幂等/四玩法/分摊/快照/流水）+ 订单查询/取消/确认 + 管理端

**Independent Test**: 幂等同单号、40001 库存不足、40006 非法迁移、秒杀并发恰好一人、取消全量回退（SC-002/003/004）

- [ ] T009 [US2] 先写失败测试：幂等（同 token 同单号/异参数冲突）、库存不足 40001、秒杀售罄 40003、拼团校验、取消回退（库存+券+积分+余额）、非法迁移 40006、分摊尾差记末行
- [ ] T010 [US2] 实现 `internal/service/shop/order_impl.go` Create：九步事务编排（data-model §一）——幂等抢占/限售/库存锁定(普通+拼团走 inventory、秒杀走 flash_sale_item 分账)/满减→券→积分→余额冻结/分摊(尾差末行)/快照落库/玩法联动(拼团成员+砍价置已下单)/状态流水
- [ ] T011 [US2] 实现 List/Detail（他人资源按不存在、状态时间线）/Cancel（条件更新+全量回退）/Confirm（→40+计提事件位）
- [ ] T012 [US2] 控制器接线：user 渠道 orders 五桩 + 转绿（含并发秒杀用例：goroutine 组抢 1 件恰 1 成功）

**Checkpoint**: 交易核心闭环（SC-002/003/004）

---

## Phase 5: User Story 3 - 支付 (Priority: P1)

**Goal**: 支付单/mock 渠道/回调幂等四层/余额联动

**Independent Test**: 发起→回调→订单 20+库存核销+余额完成；重放零变化；金额不符拒绝留档（SC-005/006）

- [ ] T013 [US3] 先写失败测试：发起（关旧待支付单）/回调成功链（订单 20+库存核销+余额完成+留档）/重放幂等/金额不符拒绝/关单
- [ ] T014 [US3] 实现 `internal/service/shop/pay_impl.go`：Create（金额=现金应付，关旧单）/Status/HandlePayNotify（四层幂等+同事务联动）/HandleRefundNotify（占位转 aftersale）/CloseExpired
- [ ] T015 [US3] 控制器接线：shop 渠道 pay 三桩 + mock 渠道报文驱动测试转绿

**Checkpoint**: 资金入口闭环（SC-005/SC-006 支付侧）

---

## Phase 6: User Story 4 - 售后 (Priority: P2)

**Goal**: 售后全链路（申请/审核/寄回/确认/退款完成联动回退）

**Independent Test**: 完成单退货退款全流程：库存回补+优惠按分摊回退+余额退回（bizType=8）+refund_status 更新

- [ ] T016 [US4] 先写失败测试：申请校验（数量/状态）→审核通过→确认收货→退款完成（mock）→库存回补+分摊回退+余额退回+refund_status；撤销；数量超额拒绝
- [ ] T017 [US4] 实现 `internal/service/shop/aftersale_impl.go`：全方法（退款 mock 渠道直调成功回调；完成时回补库存/退分摊/更新 refund_status/佣金冲销事件位）
- [ ] T018 [US4] 控制器接线：aftersale 七桩 + 转绿

**Checkpoint**: 售后资金闭环

---

## Phase 7: Polish & Cross-Cutting Concerns

- [ ] T019 quickstart 四场景终验（curl 序列）输出留证到本目录 verification.md；恒等式三条+余额式对实测数据逐条断言
- [ ] T020 宪法 IV 终验：`make build && make test` 全绿；gofmt 空输出；清单核对（SC-001~007 逐条勾验）
- [ ] T021 遗留登记：管理端权限点真实校验（治理特性）、通知投递真实发送（消息特性）、支付真实渠道——写入 tasks Notes 与 spec Assumptions 对照

---

## Dependencies & Execution Order

- Setup（T001~T003）→ Foundational（T004/T005）→ US1（T006~T008）→ US2（T009~T012）→ US3（T013~T015）→ US4（T016~T018）→ Polish（T019~T021）
- **硬依赖**：T009 下单依赖 T007 试算（共享优惠计算）；T013 支付依赖 T011 订单；T016 售后依赖 T011/T015（订单+支付）
- 并行：T001/T002/T003 之间；各故事测试先行与实现解耦

### Parallel Opportunities

- T001（money 库）与 T003（mock 渠道）与 T002（dao）三者并行
- 测试编写任务始终先行一步（红），实现任务依次转绿

---

## Implementation Strategy

### MVP First

Setup → Foundational → US1（试算）→ US2（下单）= 核心交易可演示（不含支付闭环）。

### Incremental Delivery

US2+US3 = 交易资金闭环（第一可上线增量）；US4 补售后收口；每故事一次中文提交、测试随行、独立评审。

---

## Notes

- 金额一律分域 int64（money 库），出参转元 string——AGENTS 红线
- 库存四操作与留痕同事务（ADR-0001）；秒杀走 flash_sale_item 分账
- 全部状态迁移条件 UPDATE，affected=0 → 40006
- 跨域（券/积分/余额）经 ports.go 接口回调，依赖倒置；禁 service/user 与 service/shop 互 import
- 验证以输出为证；每任务组一次提交
