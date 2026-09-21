# Contract: 008-store 端点 × 权限点映射

> 7 端点全景（spec SC-001 对账基准）。权限列延续批次 01 `middleware.RequirePerm` 显式挂接（调用点即文档）。
> 零新增错误码：资源不存在=10006、参数/区划码非法=10001、权限不足=10005（均既有）。

## admin 渠道（5）

| 端点 | 方法+路径 | 鉴权 | 权限点 |
|---|---|---|---|
| AdminStoreList | GET /admin/stores | 登录 | — |
| AdminStoreCreate | POST /admin/stores | 登录 | store:manage:create |
| AdminStoreUpdate | PUT /admin/stores/{id} | 登录 | store:manage:update |
| AdminStoreDelete | DELETE /admin/stores/{id} | 登录 | store:manage:delete |
| AdminStoreDetail | GET /admin/stores/{id} | 登录 | — |

## common 渠道（2，公开）

| 端点 | 方法+路径 | 鉴权 | 说明 |
|---|---|---|---|
| StoreList | GET /common/stores | 公开（/common/ 全公开） | 区县筛选或附近检索 |
| StoreDetail | GET /common/stores/{id} | 公开 | 歇业店可见（FR-006） |

## 错误码（全部复用既有）

| 码 | 触发 |
|---|---|
| 10001 | 区划码非 6 位数字 / 非法 ID |
| 10005 | 无权写操作（RequirePerm） |
| 10006 | 门店不存在（详情/修改/删除目标；游客详情歇业店可见不受影响） |
