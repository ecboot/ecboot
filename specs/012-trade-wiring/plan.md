# Implementation Plan: 交易闭环连线（012-trade-wiring）

**Branch**: `012-trade-wiring` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

交易闭环收口：① 前置修复积分抵扣 bug（`promotion_calc` 查询不存在的列）；② 连线购物车 5 + 订单 C 端 5；
③ 新写支付 5（`PayLogicImpl`）、后台订单 5 + 订单确认收货 1 + 定时 2（`OrderLogicImpl` 管理面）、
会员优惠券 6（`IUserCouponLogic`）。**零表结构变更、零新迁移、零新 DTO**（既有 100 个类型已覆盖）。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3

**Primary Dependencies**: 既有栈；`library/paychannel`（mock 渠道，006 已交付）；零新依赖

**Storage**: MySQL 8.4（trade_order/item/log、pay_order、pay_callback_log、coupon、user_coupon 全既有）

**Testing**: `internal/service/shop/*_impl_test.go`（复用 006 的 `trade_test.go` 基座与 fixture）、
`internal/service/user/coupon_impl_test.go`（复用 011 基座）；TDD 施于全部新写方法

**Target Platform**: `service/shop`（pay_impl 新增；order_impl 追加管理面）、`service/user`（coupon_impl 新增）、
`service/shop/promotion_calc.go`（bug 修复）；controller：shop 14 + admin 5 + user 3 桩填充

**Project Type**: GoFrame 模块化单体·交易闭环收口

**Constraints**: 分层契约（controller 只调 service）；形态跟随域内既有（shop 用 struct、user 用包级）；
**支付回调幂等四层防线为硬约束**（条件更新/唯一约束/金额校验/原文留档）

**Scale/Scope**: 22 端点、19 个新写方法 + 1 处修复、0 迁移、0 新 DTO

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落既有 shop/user 域与既有 controller；跨域（订单→券/积分/库存）经既有 ports 与事务编排 |
| II 统一技术栈 | ✅ | 零新依赖（mock 渠道复用 006 的 paychannel） |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | 新写方法 TDD；幂等/金额校验行为级断言；`check-stub` 数字 |
| V 简单优先 | ✅ | 管理面追加到既有 `OrderLogicImpl`（不新建结构）；零迁移零新 DTO |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD |
| 资金安全（既有约定） | ✅ | 金额 decimal；回调幂等四层防线；原文留档 |

**Phase 1 复查**：✅ 四产物无越界；前置修复属本批文件边界（批次 05 评审已记账）。

## Project Structure

### Documentation (this feature)

```text
specs/012-trade-wiring/{plan,research,data-model,quickstart}.md,
contracts/trade-endpoint-mapping.md, checklists/requirements.md, tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── internal/service/shop/
│   ├── promotion_calc.go        # 前置修复（1 行: 移除不存在列引用）
│   ├── pay_impl.go              # 新增: Create/Status/HandlePayNotify/HandleRefundNotify/CloseExpired
│   ├── order_impl.go            # 追加: Confirm/AdminList/AdminDetail/Deliver/AdminCancel/SellerRemark/CancelTimeout/AutoConfirm
│   └── pay_impl_test.go / order_mgmt_test.go   # 新增测试
├── internal/service/user/
│   ├── coupon_impl.go           # 新增: 6 方法（3 端点 + 3 内部）
│   └── coupon_impl_test.go
└── internal/controller/{shop,admin,user}/      # 22 桩填充
```

**Structure Decision**: 新增 3 个实现对文件（pay/coupon/mgmt 测试），其余为既有文件追加与桩填充。

## Complexity Tracking

> 无违例。管理面追加到既有 struct 而非新结构（宪法 V）。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
