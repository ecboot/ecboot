# Data Model: 009-logistics-ops

> 零表结构新增；DTO 新增 7 个、错误码 1 个。

## 一、实体

### logistics_company（物流公司字典）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| id | 主键 | 路径参数 |
| code | 编码 | **唯一（uk_code）**；创建必填、**不可修改**（订单 deliver_company 存此值） |
| name | 名称 | 创建必填 |
| tracking_rule | 运单号校验规则描述 | 可选（如 "SF+12位数字"） |
| sort | 排序 | 越小越靠前 |
| status | 1启用 0停用 | **停用保留**（列表/详情可见，仅不再供新发货选择） |
| deleted | 软删 | 1 后列表与详情均不可见 |

### operation_banner（轮播/弹窗）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| position | 1首页轮播 2首页弹窗 | 创建必填（api `in:1,2`） |
| image_url | 图片 | 创建必填 |
| link_url | 跳转链接 | 可空（空=纯展示） |
| sort | 排序 | 越小越靠前；列表主排序键 |
| start_time / end_time | 投放起止 | **可空**：起=立即、止=长期；C 端在投判定依据 |
| status | 1启用 0停用 | C 端仅返回启用 |
| deleted | 软删 | C 端与管理端均不可见 |

### operation_floor（首页楼层）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| floor_type | 1金刚区 2商品楼层 3专题 | 创建必填（api `in:1,2,3`） |
| title | 标题 | 可选 |
| config | JSON 对象 | **商品楼层约定 `{"spuIds":["<id>",...]}`（D2）**；金刚区/专题不透明，原样返回 |
| sort / status / deleted | 同 banner 语义 | C 端仅返回启用未删 |

## 二、C 端在投判定（FR-010, research D4）

```
WHERE status=1 AND deleted=0 AND position=?
  AND (start_time IS NULL OR start_time <= NOW())
  AND (end_time   IS NULL OR end_time   >= NOW())
ORDER BY sort ASC, id ASC
```

## 三、商品摘要装配（FR-012, research D3）

```
floor.config.spuIds ──▶ product_spu WHERE id IN (spuIds) AND status=1 AND deleted=0
                          ├─ name        → name
                          ├─ images[0]   → image（复用 firstImage helper）
                          └─ price_min   → price（元, 与商品列表口径一致）
失效/不存在/下架商品被逐个剔除（不报错, 数量可少于 spuIds）
```

## 四、端点 × service 方法 × 权限

| 端点 | service 方法 | 权限 |
|---|---|---|
| GET /admin/logistics-companies | LogisticsList | 登录 |
| POST /admin/logistics-companies | LogisticsCreate | logistics:company:manage |
| PUT /admin/logistics-companies/{id} | LogisticsUpdate | logistics:company:manage |
| DELETE /admin/logistics-companies/{id} | LogisticsDelete | logistics:company:manage |
| GET /admin/logistics-companies/{id} | LogisticsDetail（**接口微扩** D7） | 登录 |
| GET /admin/banners | BannerList | 登录 |
| POST /admin/banners | BannerCreate | operation:banner:manage |
| PUT /admin/banners/{id} | BannerUpdate | operation:banner:manage |
| DELETE /admin/banners/{id} | BannerDelete | operation:banner:manage |
| GET /admin/floors | FloorList | 登录 |
| POST /admin/floors | FloorCreate | operation:floor:manage |
| PUT /admin/floors/{id} | FloorUpdate | operation:floor:manage |
| DELETE /admin/floors/{id} | FloorDelete | operation:floor:manage |
| GET /shop/banners | PublicBanners | 公开（白名单已有） |
| GET /shop/floors | PublicFloors | 公开（白名单已有） |

## 五、新增 DTO（7 个, model/dto_shop.go, 域前缀命名）

- 管理面：`OperBannerItem` / `OperBannerInput` / `OperFloorItem` / `OperFloorInput`
- C 端：`PublicBannerItem`（Id/ImageUrl/LinkUrl）/ `PublicFloorItem`（FloorId/FloorType/Title/Config/Products）/
  `FloorProductSummary`（SpuId/Name/Image/Price）

## 六、错误码（新增 1）

| 码 | 语义 | 触发 |
|---|---|---|
| 40012 | 物流编码已存在 | LogisticsCreate 冲突（交易域段） |
| 10001/10006（既有） | 参数非法 / 资源不存在 | 位置与类型域外由 api `in:` 拦截；目标不存在 |
