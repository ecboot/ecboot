# 端点到实现映射：营销后台（016-marketing-admin）

**本批 30 个端点**（全为 `CodeNotImplemented` 桩；验收口径：`check-stub` admin 58 → 28）。

## 一、接口（已有，本批补实现；微扩见 D3）

```go
// ICouponLogic（service/shop/promotion.go）——补 AdminDetail（D3-①）
AdminList(ctx, status int, page model.PageReq) (*model.PageResult[model.CouponTemplate], error)
AdminDetail(ctx context.Context, id int64) (*model.CouponTemplate, error) // ← 微扩
AdminCreate(ctx, in model.CouponInput) (int64, error)
AdminUpdate(ctx, id int64, in model.CouponInput) error
AdminDelete(ctx, id int64) error
AdminRecords(ctx, couponId int64, page model.PageReq) (*model.PageResult[model.CouponRecordItem], error)

// IActivityLogic（28 方法, 全实现）: 满减5 + 拼团6 + 秒杀6 + 砍价6 + 助力5
```

## 二、端点映射（30）

| # | 端点 | controller 桩 | service 方法 | 权限 |
|---|---|---|---|---|
| 1 | `GET /admin/coupons` | admin_v1_admin_coupon_list.go | `ICouponLogic.AdminList` | coupon:read |
| 2 | `GET /admin/coupons/{id}` | admin_v1_admin_coupon_detail.go | `.AdminDetail`（D3-①） | coupon:read |
| 3 | `POST /admin/coupons` | admin_v1_admin_coupon_create.go | `.AdminCreate` | coupon:create |
| 4 | `PUT /admin/coupons/{id}` | admin_v1_admin_coupon_update.go | `.AdminUpdate` | coupon:update |
| 5 | `DELETE /admin/coupons/{id}` | admin_v1_admin_coupon_delete.go | `.AdminDelete` | coupon:delete |
| 6 | `GET /admin/coupons/{id}/records` | admin_v1_admin_coupon_record_list.go | `.AdminRecords` | coupon:read |
| 7-11 | `GET/POST /admin/full-reductions`, `GET/PUT/DELETE /admin/full-reductions/{id}` | admin_v1_admin_full_reduction_{list,create,detail,update,delete}.go | `IActivityLogic.FullReduction{List,Create,Detail,Update,Delete}` | fullreduction:manage |
| 12-16 | `GET/POST /admin/group-buys`, `PUT/DELETE /admin/group-buys/{id}`, `PUT /admin/group-buys/{id}/items` | admin_v1_admin_group_buy_{list,create,update,delete,items}.go | `GroupBuy{List,Create,Update,Delete,SetItems}` | groupbuy:manage |
| 17-21 | `GET/POST /admin/flash-sales`, `PUT/DELETE /admin/flash-sales/{id}`, `PUT /admin/flash-sales/{id}/items` | admin_v1_admin_flash_sale_{list,create,update,delete,items}.go | `FlashSale{List,Create,Update,Delete,SetItems}` | flashsale:manage |
| 22-26 | `GET/POST /admin/bargains`, `PUT/DELETE /admin/bargains/{id}`, `PUT /admin/bargains/{id}/items` | admin_v1_admin_bargain_{list,create,update,delete,items}.go | `Bargain{List,Create,Update,Delete,SetItems}` | bargain:manage |
| 27-30 | `GET/POST /admin/assists`, `PUT/DELETE /admin/assists/{id}` | admin_v1_admin_assist_{list,create,update,delete}.go | `Assist{List,Create,Update,Delete}` | assist:manage |

**无端点消费的接口方法**（D4, 记账）：`GroupBuyDetail`/`FlashSaleDetail`/`BargainDetail` 实现但不挂端点（api 无读端点; items 只写不读）。

## 三、校验矩阵（写侧）

| 对象 | 校验 | 失败码 |
|---|---|---|
| 券创建 | type=1 须 threshold；type=2 无门槛；validType=1 须起止且 start<end；validType=2 须 validDays≥1；discount>0 且 ≤2 位小数 | 10001 |
| 满减 | 档位非空、同门槛唯一、threshold/discount 合法；scope: type1 无 target、type2/3 必须有合法 targetId；时间窗 start<end | 10001 / 50008 |
| 拼团 items | groupPrice>0 合法金额；SKU 存在未删 | 10001 / 30001 |
| 秒杀 items | flashPrice>0；stockCount≥1；perLimit≥1（缺省 1） | 10001 |
| 砍价 items | original>floor>0；maxCutCount≥0（0=不限） | 10001 |
| 助力 | rewardType=1 须 rewardRef 为存在未删券；rewardType=2 须 pointAmount≥1；requiredCount≥1；时间窗 | 10001 / 50003 |
| 已售保护 | 秒杀替换移除"已被订单引用且 sold_count>0"的场次商品 → 拒绝 | 50002 |
| 撞唯一键 | 档位门槛 / 场次商品 uk_activity_sku | 50008 / 50002（1062 转业务码） |
