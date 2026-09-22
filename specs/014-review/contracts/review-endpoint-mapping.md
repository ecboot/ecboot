# 端点到实现映射：评价域（014-review）

**本批 4 个端点**（当前全为 `CodeNotImplemented` 桩；验收口径：`check-stub` 的 shop 桩数 **16 → 12**）。

| # | 端点 | controller 桩文件 | service 方法 | 鉴权 | 行为 |
|---|---|---|---|---|---|
| 1 | `POST /shop/reviews` | `shop_v1_review_create.go` | `NewReviewLogic().Create` | 会员 | 校验归属/已完成/未评过 → 落评价（audit_status=1） |
| 2 | `POST /shop/reviews/{reviewId}/extra` | `shop_v1_review_extra.go` | `.Extra` | 会员（仅本人） | 条件更新追评（一次 + 90 天内） |
| 3 | `GET /shop/reviews/mine` | `shop_v1_my_review_list.go` | `.MyList` | 会员 | 本人评价分页（含审核状态/追评/回复） |
| 4 | `GET /shop/products/{spuId}/reviews` | `shop_v1_product_review_list.go` | `.ProductList`（**本批微扩**） | 会员/公开（走既有 shop 浏览鉴权口径） | 商品评价列表 + 汇总（只出审核通过） |

**契约微扩（记账在案）**: `IReviewLogic` 原缺 `ProductList` —— 第 4 个端点在 api 契约中存在（`api/shop/v1/product.go`），service 接口漏定义。本批按批次 01/02/04 先例补：

```go
// ProductList 商品评价列表与汇总（FR-007/008）: 只出审核通过; score>0 时按星级筛选。
ProductList(ctx context.Context, spuId int64, score int, page model.PageReq) (*model.ReviewSummary, *model.PageResult[model.ReviewCard], error)
```

**返回形态**: 汇总与列表**一起算**（同源同筛），避免调用方两次查询口径漂移（research D5）。

## 需求 → 端点覆盖

| FR | 端点 |
|---|---|
| FR-001/002/003/004 | 1 |
| FR-005/010 | 3 |
| FR-006 | 2 |
| FR-007/008/009/010 | 4 |
| FR-011 | 1、2（唯一键兜底 + 条件更新） |
| FR-012 | 全部（桩清零） |

## 分层与形态约束（`docs/layer-contracts.md`）

- controller 只做"绑定 → 调 service → 映射响应"；`shop` 渠道 C 端会员端点一律先过 `requireMember`（批次 06 的 I8 防线）。
- 实现形态：**struct 方法**（`ReviewLogicImpl` + `NewReviewLogic()`），与 `CartLogicImpl`/`AfterSaleLogicImpl` 同形。
- 跨域：展示评价人需读 `user` 表的昵称——**跟随既有先例**（`order_mgmt_impl.go` 已直接读 `dao.User`），**不新增 ports**（research D2）；`maskNickname` 因兄弟域隔离而在 shop 域内自带等价实现。
- 生成物：本批**无 dao/model 变更**（表结构未动）。
- **死契约不动**: `Reply`/`Audit`/`PendingAudit` 三个后台方法无 API 端点，本批不实现、不删除，记账待后续后台评价面（research D7）。
