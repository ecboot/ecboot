# Data Model: 008-store

> 零表结构新增、零新 DTO、零新错误码。权威定义以 000031 迁移与 gf 生成物为准。

## 一、实体：store（线下门店）

| 字段 | 语义 | 本批规则 |
|---|---|---|
| id | 主键 | 路径参数、排序次键 |
| store_no | 门店编码 | **全局唯一（含已删行）**；服务端生成 `"ST"+sonyflake ID`，冲突重试 ≤3 次 |
| name | 名称 | 创建必填；关键词筛选命中字段之一 |
| province_code / city_code / district_code | 区划码 | GB/T 2260 六位数字；创建必填（服务端校验） |
| detail_address | 详细地址 | 创建必填 |
| longitude / latitude | 经纬度 | GCJ-02，DECIMAL(10,6)，**可空**；空值门店不参与附近检索 |
| business_hours / contact_phone | 营业时间/电话 | 可选档案 |
| pickup_enabled | 自提开关 | 0/1；如返回实值（交易侧消费归 012+） |
| sort | 运营排序 | 越小越靠前；区县/后台列表主排序键 |
| status | 1 营业 2 歇业 | 游客列表仅营业；游客详情歇业可见 |
| deleted | 软删 | 1 后后台与游客均不可见（资源不存在） |

## 二、状态机

```
创建 → 营业(1) ⇄ 歇业(2)      （AdminUpdate 切换）
任意态 → 软删(deleted=1)       （AdminDelete, 终态；编码不复活不复用）
```

## 三、附近检索数据流（FR-003/004, research D1）

```
入参: lng, lat, radiusKm(默认10, 上限100, 域外回退)
① 包围盒预筛(typed Where, 命中 idx_location):
   latitude  BETWEEN lat - Δlat  AND lat + Δlat     Δlat  = radiusKm / 111.0
   longitude BETWEEN lng - Δlng  AND lng + Δlng     Δlng  = radiusKm / (111.0 * cos(lat))
   AND status=1 AND deleted=0                       （游客口径）
② 计算列: Haversine(lat₁,lng₁,lat₂,lng₂) * 1000 AS distance_m   （地球半径 6371km）
③ HAVING distance_m <= radiusKm * 1000
④ ORDER BY distance_m ASC
```

区县模式与后台列表：`WHERE district_code=? AND status=1 AND deleted=0 ORDER BY sort ASC, id ASC`（后台不限 status，支持筛选）。

## 四、端点 × service 方法映射

| 端点 | service 方法 | 权限 |
|---|---|---|
| GET /admin/stores | AdminList | 登录 |
| POST /admin/stores | AdminCreate | store:manage:create |
| PUT /admin/stores/{id} | AdminUpdate | store:manage:update |
| DELETE /admin/stores/{id} | AdminDelete | store:manage:delete |
| GET /admin/stores/{id} | AdminDetail（**接口微扩**） | 登录 |
| GET /common/stores | PublicList | 公开 |
| GET /common/stores/{id} | PublicDetail | 公开 |
