# Phase 0 研究与决策：营销 C 端（015-marketing-c）

**Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

本批性质：**12 个桩 + C 端营销接口与 DTO 整体缺失（需新建）** + **清偿批次 07 的秒杀欠账**。

## D1 秒杀欠账的清偿口径（批次 07 已记账移交）

- **Decision**: ①`collectLines` 按 `flash_sale_item.flash_price` 取成交价（`price` 落秒杀价、`original_price` 落商品原价）；②下单**同时**锁活动库存（`stock_count - sold_count >= ?` 条件更新）与 `inventory`；③校验"所选 SKU 属于该活动项"（`uk_activity_sku(activity_id, sku_id)` 唯一键 → 按二元组查询）；④`cancelBy` 回补 `sold_count`（与 `inventory.locked` 一并）；⑤**删除批次 07 的临时拦截**（`FlashSaleItemId>0 → 50002`）。
- **Rationale**: 批次 07 的 ledger 明写该闸的删除条件是"接通秒杀价快照 + inventory 锁 + 取消回补 sold_count"。不接则秒杀单**永远无法支付**（支付回调的库存核销必未命中而回滚）——这是批次 07 实测过的结论。
- **Alternatives considered**: 不锁 inventory 而让支付回调跳过核销（会在回调里引入"按活动类型分支"，且商品库存永远不扣 = 账实不符，否决）。

## D2 首页聚合的组成（用户裁定）

- **Decision**: `/index` 一次返回 4 块：**轮播**（复用 `IOperationLogic.PublicBanners(1)`）、**楼层**（`PublicFloors()`）、**五类活动入口**（各取进行中前 N 条，复用本批的公开列表）、**可领券**（复用 `ICouponLogic.PublicList`）。任一块为空返回空数组。
- **Rationale**: `/index` 的 api 契约是脚手架占位（`IndexRes{Msg string}`），需设计；用户裁定"用现有能力组合，不新建表/DTO"（宪法 V）。
- **Alternatives considered**: 只做轮播+楼层（首页缺活动入口）；再加推荐商品（需新定义推荐口径，超出本批）。

## D3 助力发奖：端口投递意图（用户裁定）

- **Decision**: 达成人数 → 置 `status=2 已完成发奖` + `finish_time`，并通过**端口** `IAssistReward` 投递发奖意图（`reward_type` 1 券 / 2 积分 + `reward_ref`）；未装配时**降级告警**。
- **Rationale**: 实际发放跨域（user 域的券/积分），本批只投递意图——与批次 07 的 `ICommissionReverse` 同型（"事件出口留出口、结算后置"）。
- **Alternatives considered**: 本批直接打通发券（要动兄弟域 + 装配层，超出本批文件边界）。

## D4 砍价金额算法：确定性均分 + 末刀补齐

- **Decision**: 每刀金额 = `(original_price − floor_price) / max_cut_count` **向下取整到分**；当"已砍 + 本刀"会低于底价时，本刀改为**恰好补齐到**
底价（`current_price − floor_price`），并使状态变为"到底价"。当前价**永不低于底价**。
- **Rationale**: `bargain_item` 提供了 original/floor/max_cut_count 三要素，无需引入随机（随机金额不可测且用户体验不可预期）。末刀补齐保证多刀之和恒等于"起始价 − 底价"（消尾差，与批次 07 售后退款同思路）。
- **Alternatives considered**: 随机金额（不可测、可能超砍）；固定步长（与"最多刀数"脱节，可能砍不到底价）。

## D5 一人一刀 / 一人一助力：唯一索引兜底 + 1062 转业务码

- **Decision**: 帮砍/助力写入依赖既有唯一键 **`uk_record_helper(record_id, helper_user_id)`** 兜底；冲突（1062）转业务码；计数字段（`bargain_record.cut_count/current_price`、`assist_record.helper_count`）用**条件更新**推进。
- **Rationale**: 与批次 08 评价域同型（`uk_order_item`）；勘察确认两个 helper 表都有该唯一键 ✓。并发重复点击只能成功一次。
- **Alternatives considered**: 行锁事务（唯一索引已存在时属冗余）。

## D6 风控挂载：端口降级（迁移明确要求）

- **Decision**: 帮砍/助力入口调用**端口** `IRiskHit`（shop 定义、装配层注入 `system.IRiskLogic` 实现）；未装配时**放行 + 告警**。
- **Rationale**: 迁移 000029 两处注释明写"砍价/助力是被刷重灾区——帮砍/助力入口**须挂风控**"，属 schema 级要求；而风控实现属批次 12 → 端口是既定解法（`ICouponQuery`/`ICommissionReverse` 先例）。
- **Alternatives considered**: 本批实现风控规则引擎（越界，批次 12 域）。

## D7 C 端营销接口与 DTO 新建（先例：批次 03）

- **Decision**: 新建三个接口与实现：`IMarketingLogic`（5 个公开列表 + `Index`）、`IBargainLogic`（Launch/Progress/Cut）、`IAssistLogic`（Launch/Progress/Help）；并新增 C 端 DTO（公开活动项/秒杀项/进度视图/帮砍·助力条目/首页聚合 等）。
- **Rationale**: 勘察确认 C 端营销**无任何接口与 DTO**（`IActivityLogic` 全是管理面）；同批次 03 因装修域缺契约而新建 `IOperationLogic`。
- **Alternatives considered**: 复用 `IActivityLogic`（其方法签名是管理面语义，混用会让 C 端与管理面的有效期过滤、分页口径纠缠）。

## D8 列表与首页入口的排序/上限口径

- **Decision**: 五类公开列表统一只出"启用 + 未删 + 在时间窗内"；**按活动结束时间升序**（最快结束的先展示）；秒杀额外含"预告"（未开始但已配置）；首页各入口取**前 N=3** 条（与列表共用同一查询函数，避免口径漂移）。
- **Rationale**: 批次 08 的教训（列表与汇总口径漂移）→ 同一函数、同一筛选，入口只是加了 `Limit`。
- **Alternatives considered**: 首页自己写一套查询（口径必然漂移，否决）。

## 未知项收敛

| 未知 | 结论 |
|---|---|
| 是否新增迁移 | **否**（表/索引齐全：`flash_price`、`uk_activity_sku`、`uk_record_helper` 均在） |
| 秒杀价来源 | `flash_sale_item.flash_price`（D1） |
| 发奖/风控 | 端口（D3/D6，用户裁定） |
| 首页内容 | 轮播+楼层+活动入口+券（D2，用户裁定） |
| 新接口/DTO | 是（D7，批次 03 先例） |

**结论**：无遗留 NEEDS CLARIFICATION，可进入 Phase 1。
