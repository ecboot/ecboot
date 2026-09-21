# Contract: 010-product-admin 端点 × 既有方法 × 权限码（连线清单）

> 27 端点全景（spec SC-001 对账基准）。"来源"列注明 impl 已存在（005 遗产）或本批新写。
> controller 调用形态跟随同域既有 struct：`shop.NewProductLogic().Xxx(ctx, ...)`。

## admin 渠道（22）

### 分类（4）
| 端点 | 方法+路径 | 既有方法 | 权限点 |
|---|---|---|---|
| AdminCategoryTree | GET /admin/categories | `AdminCategoryTree` | 登录 |
| AdminCategoryCreate | POST /admin/categories | `AdminCategoryCreate` | product:category:create |
| AdminCategoryUpdate | PUT /admin/categories/{id} | `AdminCategoryUpdate` | product:category:update |
| AdminCategoryDelete | DELETE /admin/categories/{id} | `AdminCategoryDelete` | product:category:delete |

### 品牌（4）
| 端点 | 方法+路径 | 既有方法 | 权限点 |
|---|---|---|---|
| AdminBrandList | GET /admin/brands | `AdminBrandList` | 登录 |
| AdminBrandCreate | POST /admin/brands | `AdminBrandCreate` | product:brand:create |
| AdminBrandUpdate | PUT /admin/brands/{id} | `AdminBrandUpdate` | product:brand:update |
| AdminBrandDelete | DELETE /admin/brands/{id} | `AdminBrandDelete` | product:brand:delete |

### SPU（7）
| 端点 | 方法+路径 | 既有方法 | 权限点 |
|---|---|---|---|
| AdminProductList | GET /admin/products | `AdminProductList` | 登录 |
| AdminProductCreate | POST /admin/products | `AdminProductCreate`（返回 id+编码） | product:spu:create |
| AdminProductDetail | GET /admin/products/{spuId} | `AdminProductDetail` | 登录 |
| AdminProductUpdate | PUT /admin/products/{spuId} | `AdminProductUpdate` | product:spu:update |
| AdminProductDelete | DELETE /admin/products/{spuId} | `AdminProductDelete` | product:spu:delete |
| AdminProductStatus | POST /admin/products/{spuId}/status | `AdminProductStatus` | product:spu:update |
| AdminProductRestrict | PUT /admin/products/{spuId}/restrict-codes | `AdminProductRestrict` | product:spu:update |

### SKU（4）
| 端点 | 方法+路径 | 既有方法 | 权限点 |
|---|---|---|---|
| AdminSkuCreate | POST /admin/products/{spuId}/skus | `AdminSkuCreate` | product:sku:create |
| AdminSkuUpdate | PUT /admin/skus/{skuId} | `AdminSkuUpdate` | product:sku:update |
| AdminSkuDelete | DELETE /admin/skus/{skuId} | `AdminSkuDelete` | product:sku:delete |
| AdminSkuStatus | POST /admin/skus/{skuId}/status | `AdminSkuStatus` | product:sku:update |

### 库存（3，**本批新写**）
| 端点 | 方法+路径 | 新写方法 | 权限点 |
|---|---|---|---|
| AdminInventoryList | GET /admin/inventories | `InventoryList` | 登录 |
| AdminInventoryWarnings | GET /admin/inventories/warnings | `InventoryWarnings` | 登录 |
| AdminInventoryAdjust | POST /admin/inventories/{skuId}/adjust | `InventoryAdjust` | inventory:adjust |

## shop 渠道（5，公开——白名单已有 `/shop/categories`、`/shop/brands`、`/shop/products`、`/shop/search`）

| 端点 | 方法+路径 | 既有方法 |
|---|---|---|
| CategoryTree | GET /shop/categories | `CategoryTree` |
| BrandList | GET /shop/brands | `Brands` |
| ProductList | GET /shop/products | `Products` |
| ProductDetail | GET /shop/products/{spuId} | `ProductDetail`（viewerUserId 从 ctx 取，未登录=0） |
| ProductSearch | GET /shop/search | `Search` |

## 库存语义（本批新写部分, 详见 data-model）

- 可售 = total − locked（推导, 不落库）
- 预警：可售 ≤ warn_count
- 调整：delta 有符号；`UPDATE inventory SET total=total+delta WHERE sku_id=? AND total+delta>=0`
  （affected=0 → 30007 库存调整非法）；同事务写入 `inventory_log`（变更前后快照 + operator + remark）
- operator：**表注释约定格式 `admin:{id}`**（inventory_log.operator dc:"操作者:system/user:{id}/admin:{id}"），
  取 `middleware.CtxUserIdFrom(ctx)` 的登录后台账号 ID 拼装
- change_type = 5（后台调整）；quantity 为绝对值；留痕含 total_after/locked_after 快照（对账用）
