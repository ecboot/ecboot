# Implementation Plan: 交易链路实现（trade-loop）

**Branch**: `006-trade-loop` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

实现核心交易域：购物车试算、订单创建（四玩法幂等）、支付（mock 渠道+回调幂等四层）、售后全链路（退款联动库存/优惠/余额/佣金事件位）。落位 `internal/service/shop` 的 cart/order/pay/aftersale 接口实现 + `internal/library/paychannel`（mock 渠道适配）。零新表，全部复用既有 68 表结构。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（gdb 事务 + dao 生成物 + do 强类型）

**Primary Dependencies**: 既有栈全部复用；零新第三方依赖

**Storage**: MySQL 8.4（事务：`g.DB().Begin()` 传递 `*gdb.TX` 至各域锁定/核销方法——inventory.Lock 等的 `tx any` 参数即此）；Redis（幂等 token、秒杀预热可选）

**Testing**: service 层 gtest 红绿（确定性配置注入 + 数据自建清理——005 模式）；并发秒杀用 goroutine 组 + `sync.WaitGroup`（`wg.Go`，modern Go）

**Target Platform**: `internal/service/shop`（cart/order/pay/aftersale 实现）、`internal/library/paychannel`（mock）

**Project Type**: 核心交易域实现

**Performance Goals**: 试算/下单 p95 < 300ms（本地容器）

**Constraints**: 分层契约 §二（do/entity 强类型禁 map）；金额内部分（int64）计算；恒等式 SC 强约束；TDD 红绿；管理端权限点挂接（延续 005 形态）

**Scale/Scope**: `service/shop` 四个接口实现（约 35 方法）、paychannel 组件、零新表、零新端点（004 契约已定）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位 service/shop 既有接口；跨域（user 积分/余额/券）经参数传入的回调接口协作（接口在 shop 定义，user 侧注册实现——依赖倒置，不 import） |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 注释/文档中文 |
| IV 可验证交付 | ✅ | 红绿 TDD + 恒等式断言 + quickstart 序列 |
| V 简单优先 | ✅ | 事务用 gdb 原生 Begin/Commit（不引事务管理器）；mock 渠道直调回调 |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD；B2C 定位 |
| 跨域协作 | ✅ | 券核销/退回、积分消耗/回退、余额冻结/解冻——经 `user` 域已注册的接口回调（依赖倒置），接口签名在 plan §跨域接口锁定 |

**Phase 1 复查**：✅ 四产物无越界；跨域接口签名已锁定（见 data-model §三）。

## Project Structure

### Documentation (this feature)

```text
specs/006-trade-loop/{plan,research,data-model,quickstart}.md, contracts/tx-invariants.md, tasks.md(后续)
```

### Source Code (apps/server)

```text
internal/service/shop/
├── cart_impl.go        # 试算（复用 005 Checkout 扩展优惠计算）
├── order_impl.go       # Create/List/Detail/Cancel/Confirm/Admin*
├── pay_impl.go         # 发起/状态/回调/关单
├── aftersale_impl.go   # 售后全链路
└── trade_test.go       # 交易域集成测试（红绿）
internal/library/paychannel/
└── mock.go             # mock 渠道适配（生成唤起参数/模拟回调报文）
```

**Structure Decision**: 订单编排集中在 `order_impl.go` 的 `Create`（唯一事务入口）；库存/券/积分/余额各域暴露**事务感知方法**（接受 `*gdb.TX`），由订单事务统一提交/回滚——单事务保证一致性，避免Saga复杂度（V1 单体最优解）。

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 无新包层级；跨域依赖倒置（接口在 shop、实现在 user 注册） |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 全中文 |
| IV 可验证交付 | ✅ | 恒等式/幂等/状态机断言 + quickstart |
| V 简单优先 | ✅ | 单事务编排；mock 渠道直调回调（不搭消息队列） |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD 红绿；B2C 定位 |
| B2C+多门店 | ✅ | 单店交易；无商户语义 |

**Phase 1 复查**：✅ 通过（跨域接口锁定见 data-model §三）。

## Project Structure

### Documentation (this feature)

```text
specs/006-trade-loop/{plan,research,data-model,quickstart}.md, contracts/tx-invariants.md, tasks.md(后续)
```

### Source Code (apps/server)

```text
internal/service/shop/
├── cart_impl.go        # 结算试算（AmountBook 全分项）
├── order_impl.go       # Create（幂等+四玩法+分摊）/查询/取消/确认/管理端
├── pay_impl.go         # 发起/状态/回调/关单
├── aftersale_impl.go   # 申请/审核/寄回/确认/退款重试
└── trade_test.go       # 红绿集成测试
internal/library/paychannel/mock.go   # mock 渠道适配
```

**Structure Decision**: 订单编排集中在 `order_impl.go` 的 `Create`（唯一事务入口）；库存/券/积分/余额各域暴露**事务感知方法**（接受 `*gdb.TX`），由订单事务统一提交/回滚——单事务保证一致性，避免Saga复杂度（V1 单体最优解）。

## Complexity Tracking

> 无宪法违例。V1 采用单事务编排（跨域经接口回调），是单体最优解而非妥协；微服务化时该事务边界即拆分断点（已在 ADR-0001 精神内）。
