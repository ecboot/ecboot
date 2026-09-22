# Implementation Plan: 营销 C 端（015-marketing-c）

**Branch**: `015-marketing-c` | **Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

**Note**: 本批为批次 09（`specs/PROGRESS.md` §三）：**12 个桩 + C 端营销契约整体缺失（需新建）+ 清偿批次 07 的秒杀欠账**。

## Summary

补齐营销 C 端并收官 shop 渠道：新建 `IMarketingLogic`（5 公开列表 + 首页聚合）、`IBargainLogic`（发起/进度/帮砍）、
`IAssistLogic`（发起/进度/助力）三个接口与实现；新增约 10 个 C 端 DTO；连线 12 个端点；
**跨批清偿秒杀欠账**（秒杀价快照 + 活动/商品**双库存**锁 + 取消回补 `sold_count` + 删除批次 07 的临时拦截）；
新增两个端口（`IRiskHit` 风控、`IAssistReward` 发奖意图）。原假设"零迁移"; 实际落地 000040~000042（订单头秒杀列经评审修复轮撤销改行级归属; `bargain_record.order_no` 改可空——详见 PROGRESS §五）。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（以 `go.mod` 为准）

**Primary Dependencies**: 既有栈；`library/money`（金额分运算）、`gconv`（JSON 配置与规格）；**零新依赖**

**Storage**: MySQL 8.4（`flash_sale_*`/`group_buy_*`/`bargain_*`/`assist_*`/`promotion_*` 全既有；
关键索引：`uk_activity_sku(activity_id,sku_id)`、`uk_record_helper(record_id,helper_user_id)`）

**Testing**: `internal/service/shop/{marketing_impl_test.go,play_impl_test.go}`（TDD，真实库 + 自清 fixture，
复用 `aftersale/review` 的订单/商品/会员 fixture 构造）；`internal/routes/` 补**端点级可达性测试**

**Target Platform**: `service/shop`（新增 3 接口 + 3 实现文件）；controller：shop 12 桩填充（**shop 渠道本批收官**）；
`api/shop/v1/index.go` 契约扩写

**Project Type**: GoFrame 模块化单体·领域实现补齐 + 债务清偿

**Constraints**: 分层契约；形态跟随 shop 域（struct 方法 + `NewXxxLogic()`）；
**写入一律"唯一索引兜底 + 条件更新 + 判 RowsAffected"**（012/013/014 三轮评审确立）；
列表口径"只出有效活动 + 按结束时间升序"，首页入口与列表**共用同一查询**（避免口径漂移）

**Scale/Scope**: 12 端点、3 新接口/实现、~10 新 DTO、2 新端口、1 处跨批改动（order_impl 秒杀分支与取消回补）、0 迁移

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落既有 shop 域；跨域（风控/发奖/昵称）一律走端口或既有先例；不 import 兄弟域 |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | TDD 覆盖 12 端点 + 并发/边界 + 秒杀闭环（含支付回调）；`check-stub` shop 12→0；`lint` 保持 0；端点级可达性测试 |
| V 简单优先 | ✅ | 首页用现有能力组合（不新建表/DTO）；砍价用确定性算法（不引随机）；唯一索引兜底（不上行锁事务） |
| 工程约束 §二/五 | ✅ | 零生成物变更；TDD |
| 资金安全 | ✅ | 秒杀双库存与回补均条件更新 + 判行数（批次 07 教训）；金额分运算 |

**Phase 1 复查**：✅ 四产物无越界。**偏移记账**：①`api/shop/v1/index.go` 契约由占位扩为聚合结构（用户裁定首页内容）；
②跨批改 `order_impl.go`（秒杀分支/取消回补，批次 07 ledger 明确移交）；③新增 2 端口与 ~10 DTO。均已在 `PROGRESS §五` 记账。

## Project Structure

### Documentation (this feature)

```text
specs/015-marketing-c/{plan,research,data-model,quickstart}.md,
contracts/marketing-c-endpoint-mapping.md, checklists/requirements.md, tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── api/shop/v1/index.go              # IndexRes 占位 → 聚合结构
├── internal/service/shop/
│   ├── marketing.go / marketing_impl.go        # 新建: IMarketingLogic（5 列表 + 首页）
│   ├── bargain.go   / bargain_impl.go          # 新建: IBargainLogic（发起/进度/帮砍）
│   ├── assist.go    / assist_impl.go           # 新建: IAssistLogic（发起/进度/助力）
│   ├── ports.go                                # +IRiskHit / +IAssistReward
│   ├── order_impl.go                           # 跨批: 秒杀价快照/双库存/校验/解除拦截/取消回补
│   └── {marketing,play}_impl_test.go           # 新增测试
├── internal/model/dto_shop.go                  # +~10 个 C 端 DTO
└── internal/controller/shop/shop_v1_*.go       # 12 桩填充（shop 渠道收官）
```

**Structure Decision**: 3 接口 + 3 实现 + 2 测试文件；其余为既有文件追加与桩填充；无新包、无新迁移、无生成物变更、无新表。

## Complexity Tracking

> 无违例。跨批改动与契约扩写均为**已记账的批次移交/用户裁定**，非静默扩范围。
