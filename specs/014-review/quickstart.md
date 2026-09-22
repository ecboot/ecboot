# 验证指南：评价域（014-review）

**Prerequisites**: 应用库 `ecboot`（`127.0.0.1:3306`）在迁移版本 **38**（本批零迁移）。命令在 `apps/server/` 下执行（本机无 `make`）。

## 一、验收命令（DoD 数字）

```bash
go build ./...
go test ./... -count=1                       # 全量回归（含 01~07 既有链路零退化）
golangci-lint run                            # 期望 0 issues（自批次 06 起为硬线）
grep -l CodeNotImplemented internal/controller/shop/*.go | wc -l   # 期望 12（本批 4 端点清零）
```

## 二、端到端场景（服务层测试覆盖，按序即一轮完整演示）

场景脚本落 `internal/service/shop/review_impl_test.go`（TDD 先红后绿），打真实库并自清。

1. **提交评价（US1）**
   造已完成订单 + 订单项 → `Create`（5 星/带图/非匿名）→ 断言落库：`audit_status=1`、`spu_name`/`sku_specs` 为**下单时快照**、`images` 为数组
   → 同一订单项再 `Create` → 40010（已评价），且**活库仍只有 1 条**。
2. **重复/并发提交（SC-003）**
   8 并发 `Create` 同一订单项 → 恰好 1 次成功、其余 40010；活库该 `order_item_id` 仅 1 行（`uk_order_item` 兜底 + 1062 转业务码）。
3. **非可评场景（US1）**
   订单项非本人 → 按不存在；订单未完成（待发货 20）→ 40010；评分 0/6 → 10001。
4. **商品评价列表与汇总（US2）**
   造 2 条通过 + 1 条待审 + 1 条驳回（同 SPU）→ `ProductList`（不筛星）→ 只返回 2 条；`score` 筛选生效；汇总 Total=2、分布只含那 2 条、Avg 正确 → 新提交一条（V1 自动通过）后 Total 变 3（**D1 口径的直接后果**）。
5. **匿名与脱敏（US2/FR-009）**
   匿名评价在列表中显示"匿名用户"；非匿名显示脱敏昵称（首字符 + 掩码）；空昵称给占位。
6. **我的评价（US3）**
   本人 2 条（不同审核状态）+ 他人 1 条 → `MyList` 只返回本人的 2 条且状态正确；含追评/回复的条目一并透出。
7. **追评（US4）**
   30 天前主评 → `Extra` 成功且 `extra_time` 落库 → 再次 `Extra` → 40011 → 100 天前主评 → 40011（超期）→ 他人评价 → 不存在 → 追评内容为空 → 10001。
8. **回归（SC-004）**
   批次 01~07 全部测试保持绿（本批只读订单/订单项，不写交易域）。

## 三、不做的事（边界）

- 不实现后台审核面（`Reply`/`Audit`/`PendingAudit` 无端点，属死契约）——**合规债务已记账：公开上线前必修**。
- 不改 `product_review` 表结构（零迁移）；不新增 DTO（`ReviewCreateInput`/`MyReviewCard`/`ReviewSummary` 等已存在）。
- 不改支付/订单/售后域既有实现。
