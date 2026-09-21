# Contract: 009-logistics-ops 端点 × 权限点映射

> 15 端点全景（spec SC-001 对账基准）。权限列延续批次 01 `middleware.RequirePerm` 显式挂接。
> 错误码：新增 40012（物流编码已存在）；其余复用 10001/10005/10006。

## admin 渠道（13）

| 端点 | 方法+路径 | 鉴权 | 权限点 |
|---|---|---|---|
| AdminLogisticsList | GET /admin/logistics-companies | 登录 | — |
| AdminLogisticsCreate | POST /admin/logistics-companies | 登录 | logistics:company:manage |
| AdminLogisticsUpdate | PUT /admin/logistics-companies/{id} | 登录 | logistics:company:manage |
| AdminLogisticsDelete | DELETE /admin/logistics-companies/{id} | 登录 | logistics:company:manage |
| AdminLogisticsDetail | GET /admin/logistics-companies/{id} | 登录 | — |
| AdminBannerList | GET /admin/banners | 登录 | — |
| AdminBannerCreate | POST /admin/banners | 登录 | operation:banner:manage |
| AdminBannerUpdate | PUT /admin/banners/{id} | 登录 | operation:banner:manage |
| AdminBannerDelete | DELETE /admin/banners/{id} | 登录 | operation:banner:manage |
| AdminFloorList | GET /admin/floors | 登录 | — |
| AdminFloorCreate | POST /admin/floors | 登录 | operation:floor:manage |
| AdminFloorUpdate | PUT /admin/floors/{id} | 登录 | operation:floor:manage |
| AdminFloorDelete | DELETE /admin/floors/{id} | 登录 | operation:floor:manage |

## shop 渠道（2，公开——Auth 白名单既有 `/shop/banners`、`/shop/floors`）

| 端点 | 方法+路径 | 鉴权 | 说明 |
|---|---|---|---|
| BannerList | GET /shop/banners?position= | 公开 | **在投过滤**（启用 + 时段内） |
| FloorList | GET /shop/floors | 公开 | 启用楼层 + **商品楼层摘要装配**（失效剔除） |

## 楼层 config 约定（本批定义）

| 楼层类型 | config 形态 | C 端处理 |
|---|---|---|
| 1 金刚区 | `{"entries":[{"icon":"","link":"","text":""}]}`（前端约定, 本批不校验） | config 原样返回 |
| 2 商品楼层 | `{"spuIds":["<spuId>",...]}`（**本批约定, 字符串 ID 数组**） | 装配商品摘要, 失效剔除 |
| 3 专题 | 不透明 JSON | config 原样返回 |
