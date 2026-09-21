# Research: 009-logistics-ops 物流与运营装修

> Phase 0 产出。基于代码勘察（commit e1792a6 时点）。

## D1 装修域接口与 DTO 新建（非微扩）

**结论**：新建 `IOperationLogic` 落 `internal/service/shop/misc.go`（与 ILogisticsLogic 同文件同域）；
DTO 落 `model/dto_shop.go`，带域前缀命名：
`OperBannerItem/OperBannerInput/OperFloorItem/OperFloorInput`（管理面）+
`PublicBannerItem/PublicFloorItem/FloorProductSummary`（C 端）。

**理由**：勘察确认装修域**无任何既有接口与 DTO**（misc.go 仅 ILogisticsLogic/IStoreLogic/IRiskLogic；
grep Banner/Floor 于 model 无果）。落 shop 域因：装修是商城内容（C 端首页消费）、物流是交易字典
（发货消费）——同域内聚，避免为"运营"另立顶层域（宪法 V）。管理端 admin 渠道消费同一域（如 007 的
system 域被 admin+common 消费，渠道与域解耦）。

**替代方案**：新建 `service/operation` 域——仅 2 个子域（banner/floor）不足以支撑独立域，且与
"平台级能力"（system）语义重叠；放弃。

## D2 楼层 config 约定（本批定义, 写进契约）

**结论**：商品楼层 config 约定 `{"spuIds": ["<spuId>", ...]}`（**字符串 ID 数组**，与 api 层 ID 一律
string 的防精度口径一致）；金刚区/专题 config 为不透明 JSON，原样返回不校验。

**理由**：000026 表注释明示"类型化 schema 由应用层约定"；004 契约将 config 定义为不透明
`map[string]any`。字符串数组与项目 ID 契约（int64 对外 string）一致，避免前端精度问题。

## D3 商品摘要装配（C 端楼层）

**结论**：`spuIds` → `product_spu`（**status=1 上架 且 deleted=0**）→ 摘要 {spuId, name, image, price}：
- image 复用既有 `firstImage(images JSON)` helper（browse_impl.go:90，取 images[0]）
- price 取 `price_min`（元，与既有商品列表价格口径一致，browse_impl 同源）
- **失效商品逐个剔除**（数量对不上不报错）——装配是"尽力而为"的展示语义

**理由**：同包 helper 复用避免重复实现；价格口径与商品列表一致避免同一商品两处两价。

## D4 投放时段判定（SQL 侧）

**结论**：在投条件 = `status=1 AND deleted=0 AND position=? AND (start_time IS NULL OR start_time<=NOW())
AND (end_time IS NULL OR end_time>=NOW())`，排序 `sort ASC, id ASC`。

**理由**：时段是投放的核心语义，判定放库侧（避免拉全量再过滤）；NULL 双端 = 长期有效（表注释
"NULL=立即生效/长期有效"）。起止反置不校验（spec 边界：运营自行负责，区间判断自然不命中）。

## D5 物流公司编码不可改 + 错误码

**结论**：`Update` 入参不含编码；编码重复新增错误码 `40012 物流编码已存在`（交易域段）。

**理由**：表注释"编码(唯一,订单deliver_company存此值)"——改码会割裂历史订单的承运商标识。
错误码选 4xxxx（交易支撑）而非复用 80004（后台治理·角色编码，语义窄）。

## D6 权限与公开性

**结论**：三类写操作挂 `logistics:company:manage` / `operation:banner:manage` / `operation:floor:manage`
（000032 种子 430/350/370）；admin 查询仅登录；shop 两端点公开（Auth 白名单已有 `/shop/banners`、
`/shop/floors`——批次 01 已配置）。

## D7 ILogisticsLogic 微扩 Detail（D6 模式延续）

**结论**：接口增 `Detail(ctx, id int64) (*model.LogisticsCompany, error)`；DTO 复用。

## 勘察结论（非决策）

- **物流 DTO 已就位**：`model.LogisticsCompany` + `LogisticsCompanyInput`（dto_shop.go:129/137）。
- **表齐备**：logistics_company（000015）、operation_banner/operation_floor（000026）；零新迁移。
- **装修表状态**：banner/floor 均 `status 1启用 0停用` + deleted 软删；floor.config 为 JSON 可空。
- **测试基座**：批次 01/02 模式延续，测试落 `internal/service/shop/operation_impl_test.go`。
