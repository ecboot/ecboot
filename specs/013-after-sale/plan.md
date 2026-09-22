# Implementation Plan: 售后域（013-after-sale）

**Branch**: `013-after-sale` | **Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

**Note**: 本批为批次 07（`specs/PROGRESS.md` §三）：**接口与状态机已定义、impl 待补 + 11 个 controller 桩连线**。

## Summary

把售后域从"有接口无实现"补成可用链路：新写 `AfterSaleLogicImpl`（11 个方法，含状态机与退款编排）、
连线 shop 5 + admin 6 = **11 个端点**、**唯一迁移 000038**（给 `after_sale_order` 补 `operator_id`/`fail_reason`
两列，解决"接口传 operator 而表无列可承载"的审计缺口）。退款渠道复用既有 `paychannel.Channel.Refund`
（幂等键 `out_refund_no = after_sale_no`），回调推进复用批次 06 已修正的 `HandleRefundNotify`。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（版本以 `go.mod` 为准）

**Primary Dependencies**: 既有栈；`library/paychannel`（mock 渠道，含 `Refund`）、`library/money`（分/元换算）、
`sonyflake`（售后单号，与订单号生成同风格）；**零新依赖**

**Storage**: MySQL 8.4（`after_sale_order` 既有 + 000038 两列；读 `trade_order_item`/`trade_order`，写 `inventory`）

**Testing**: `internal/service/shop/aftersale_impl_test.go`（TDD，真实库 + 自清 fixture，复用 `order_create_test.go`
的订单/地址/商品 fixture 与 `trade_test.go` 基座）；controller 层复用 `internal/controller/shop/*_test.go` 基座

**Target Platform**: `service/shop`（新增 `aftersale_impl.go`；`aftersale.go` 的接口注释按 D4 修正）；
controller：shop 5 + admin 6 桩填充；`migrations/000038_*`；`gf gen dao` 重生成 `after_sale_order` 生成物

**Project Type**: GoFrame 模块化单体·领域实现补齐

**Constraints**: 分层契约（controller 只调 service）；形态跟随 shop 域（struct 方法 + `NewXxxLogic()`）；
**全部状态迁移必须"条件更新 + 判 RowsAffected"**（批次 06 修正后的既有惯例，I1/I11 同型缺陷不得复发）；
金额一律 `money` 库 + 分运算；`refund_amount ≤ 订单项 pay_amount`

**Scale/Scope**: 11 端点、11 个新写方法、1 迁移（2 列）、0 新表、0 新 DTO（`model.AfterSale*` 三个类型已存在）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落既有 `service/shop` 域与既有 controller；读订单项/写库存均在域内，佣金冲销只投事件 |
| II 统一技术栈 | ✅ | 零新依赖（渠道/金额/编号全复用既有库） |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | TDD 覆盖状态机每一跳 + 幂等/边界；`check-stub` 数字（shop 21→16、admin 64→58）；`lint` 保持 0 issues |
| V 简单优先 | ✅ | 复用既有 `paychannel.Refund` 与 `HandleRefundNotify`，不新建渠道抽象/不新建日志表（只加 2 列） |
| 工程约束 §二/五 | ✅ | 生成物 `gf gen dao` 重生成（勿手改）；`do/entity` 强类型 |
| 资金安全（既有约定） | ✅ | 金额分运算；退款幂等键=售后单号；条件更新防重入；0 元不调渠道 |

**Phase 1 复查**：✅ 四产物无越界。**唯一偏移**：spec 原写"不新增迁移"，Phase 0 的 D3 以实证（接口传 operator
而表无承载列）推翻了该假设 → 新增 000038；已在 `research.md` D3 记录原因、代价与先例（批次 01/05 同样处理），
并须在 PROGRESS §五 记账。

## Project Structure

### Documentation (this feature)

```text
specs/013-after-sale/{plan,research,data-model,quickstart}.md,
contracts/after-sale-endpoint-mapping.md, checklists/requirements.md, tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── migrations/000038_after_sale_operator.{up,down}.sql   # 新增: operator_id + fail_reason
├── internal/model/{entity,do}/after_sale_order.go        # gf gen dao 重生成（勿手改）
├── internal/service/shop/
│   ├── aftersale.go              # 接口注释按 D4 修正（撤销边界口径）
│   ├── aftersale_impl.go         # 新增: 11 方法（Apply/List/Detail/Cancel/SubmitReturn + Admin List/Detail/Approve/Reject/ConfirmReceipt/RetryRefund）
│   └── aftersale_impl_test.go    # 新增: 状态机每一跳 + 幂等 + 边界 TDD
└── internal/controller/{shop,admin}/                      # 11 桩填充
    ├── shop_v1_after_sale_{create,list,detail,cancel,logistics}.go
    └── admin_v1_admin_after_sale_{list,detail,approve,reject,confirm_receipt,retry_refund}.go
```

**Structure Decision**: 一个实现文件（`aftersale_impl.go`）+ 一个测试文件，其余为既有文件追加与桩填充；
不新增包、不新增 DTO、不新增渠道抽象。

## Complexity Tracking

> 仅一项需说明：本批**非零迁移**（原 spec 假设零迁移）。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| 新增迁移 000038（2 列） | `IAfterSaleLogic` 的 4 个后台方法都接收 `operator`，而 `after_sale_order` 无列可承载；`admin_operation_log` 至今无人写入（属批次 12），无法作为载体。不补 = 审计静默丢失（批次 06 评审同型 Important 缺陷） | ① 不补、记缺口 → 重蹈"operator 被丢弃"的覆辙；② 复用业务列（reject_reason/description）→ 语义污染；③ 新建 after_sale_log 表 → 本批只需"最后操作人 + 失败原因"，一表过重 |
