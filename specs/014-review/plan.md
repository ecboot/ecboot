# Implementation Plan: 评价域（014-review）

**Branch**: `014-review` | **Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

**Note**: 本批为批次 08（`specs/PROGRESS.md` §三）：**接口已有三方法（Create/Extra/MyList）、一个端点缺方法、4 个 controller 桩连线**。

## Summary

补齐评价域：新写 `ReviewLogicImpl`（4 复用接口方法 + 1 微扩 `ProductList` = 5 方法）、连线 shop 4 个端点、
**零迁移零新 DTO**。三项口径裁定：① **V1 自动通过**（用户裁定 B，因无后台审核端点，否则评价列表恒空）；
② 评价人昵称**跟随既有先例**由 shop 域直读 `user` 表（包级 import 隔离仍守）+ 域内自带脱敏；
③ "一项一评"依赖既有 `uk_order_item` 唯一索引兜底，**1062 必须转业务码**；追评用**条件更新**（一次 + 90 天内）。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（以 `go.mod` 为准）

**Primary Dependencies**: 既有栈；`gconv`（JSON 数组解析）；**零新依赖**

**Storage**: MySQL 8.4（`product_review` 既有 + `uk_order_item` 唯一索引；读 `trade_order`/`trade_order_item`/`user`，**不写**交易域）

**Testing**: `internal/service/shop/review_impl_test.go`（TDD，真实库 + 自清 fixture，复用 `aftersale_impl_test.go` 的订单/订单项 fixture 构造）

**Target Platform**: `service/shop`（新增 `review_impl.go`；`review.go` 接口微扩 ProductList）；controller：shop 4 桩填充

**Project Type**: GoFrame 模块化单体·领域实现补齐

**Constraints**: 分层契约；形态跟随 shop 域（struct 方法 + `NewReviewLogic()`）；
**写入一律"唯一索引兜底 + 条件更新"**（批次 07 C1/C2 教训）；DB 唯一键冲突转业务码；列表与汇总**同源同筛**

**Scale/Scope**: 4 端点、5 个方法（4 复用 + 1 微扩）、0 迁移、0 新 DTO、0 生成物变更

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落既有 shop 域与既有 controller；读 `user` 表**跟随既有先例**（order_mgmt 已如此），包级 import 隔离仍守 |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | TDD 覆盖 4 端点 + 并发/边界；`check-stub` 16→12；`lint` 保持 0；两连跑全绿 |
| V 简单优先 | ✅ | 零迁移、零新 DTO、零新端口（昵称直读既有表）；不为一个展示字段新增装配 |
| 工程约束 §二/五 | ✅ | 无生成物变更；TDD |
| 资金安全 | N/A | 本域不涉金额 |

**Phase 1 复查**：✅ 四产物无越界。**唯一偏移**：`IReviewLogic` 补 `ProductList`（接口漏定义，同批次 01/02/04 先例），
已在 `contracts/` 与 `PROGRESS §五` 记账；**合规债务**（UGC 未审核即公开）同处记账。

## Project Structure

### Documentation (this feature)

```text
specs/014-review/{plan,research,data-model,quickstart}.md,
contracts/review-endpoint-mapping.md, checklists/requirements.md, tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── internal/service/shop/
│   ├── review.go            # 接口微扩: +ProductList（其余三方法签名不动）
│   ├── review_impl.go       # 新增: Create/Extra/MyList/ProductList + 脱敏/汇总 helper
│   └── review_impl_test.go  # 新增: 4 端点 + 并发/边界 TDD
└── internal/controller/shop/
    ├── shop_v1_review_create.go / shop_v1_review_extra.go        # 桩填充
    └── shop_v1_my_review_list.go / shop_v1_product_review_list.go # 桩填充
```

**Structure Decision**: 一个实现文件 + 一个测试文件；不新增包、不新增 DTO、不新增迁移、无生成物变更。

## Complexity Tracking

> 无违例。轻量实现（宪法 V）。
