# Research: 008-store 门店域

> Phase 0 产出。全部基于代码现状勘察（commit b3212f0 时点）。

## D1 附近检索 = 包围盒预筛 + Haversine 精确过滤（dao 链内完成）

**结论**：`PublicList` 附近模式在 `dao.Store` 链内完成——① 先以半径换算经纬度包围盒（typed Where，
命中 `idx_location`）预筛；② SELECT 追加 Haversine 球面距离表达式为计算列 `distance_m`；
③ `HAVING distance_m <= 半径*1000` 过滤 + `ORDER BY distance_m ASC`。无坐标门店被包围盒条件天然排除。

**理由**：库内数据规模（数百店）无性能专项；包围盒保证不走全表扫描；表达式仅是 dao 链上的计算列
（与批次 01 已评审通过的 Where 片段同性质），不引入裸 `g.DB()`。repository 泛型脚手架未被任何代码接线，
为单查询启用一层违反宪法 V（简单优先）；放弃 repository 具名方法方案。

## D2 门店编码 = sonyflake 首次接线（项目约定落地）

**结论**：新增 `internal/library/idgen`（library 既有层的新叶子包, 与 captcha/sms 同级）承载 sonyflake 单例；
`store_no = "ST" + sonyflake 十进制 ID`；唯一键冲突时重试（调用方无感知, 上限 3 次）。

**理由**：AGENTS.md/layer-contracts 基础库约定"分布式 ID 统一 sonyflake（业务编号场景）"——该约定
（c482eb0）晚于交易链路（40da88c），`nextOrderNo` 的"前缀+时间戳+随机数"系约定落地前产物，
**不在本批改动**（防偏离条款 3：批次边界）；门店编码是新业务编号, 适用约定。idgen 为后续订单号
向约定迁移提供共享落点。

## D3 AdminDetail 接口微扩（D6 模式延续）

**结论**：`IStoreLogic` 增加 `AdminDetail(ctx, storeId int64) (*model.StoreItem, error)`；DTO 零新增
（StoreItem 已含编码/坐标全字段）。游客详情复用 `PublicDetail`。

## D4 权限挂接延续批次 01 RequirePerm 机制

**结论**：create/update/delete 显式挂 `store:manage:create/update/delete`（000032 种子 400/390/410）；
列表/详情仅要求登录（004 契约注释只标注写操作权限点）。零新增错误码：不存在复用 10006，
区划码非法走 10001（api 层 `v:"length:6"` + 服务端 6 位数字校验双保险）。

## D5 排序口径

**结论**：区县筛选/后台列表按 `sort ASC, id ASC`（store 表自带运营排序位）；附近模式按距离 ASC。
接口契约未标注排序——以表设计（sort 注释"越小越靠前"）为准绳。

## 勘察结论（非决策）

- **DTO 零新增**：StoreQuery/StoreItem/StoreInput 已就位（dto_shop.go），PublicDetail 返回 StoreItem。
- **表结构**：store_no 唯一键（含已删行——软删编码不复活）；status 1营业/2歇业；GCJ-02 坐标可空。
- **歇业口径**：游客列表仅营业（接口契约）；游客详情歇业可见（表语义"展示但不可自提/核销"，
  交易侧裁决归 012+）。
- **测试基座**：批次 01 模式延续——但门店属 shop 域（service/shop 包），测试落 `internal/service/shop/store_impl_test.go`。
