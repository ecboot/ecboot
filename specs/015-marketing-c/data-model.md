# Phase 1 数据模型：营销 C 端（015-marketing-c）

**依据**: 迁移 V17(拼团)/V18(秒杀)/V20(满减)/V29(砍价·助力)。原假设"零迁移"; 实际经过两轮更正（详见 PROGRESS §五）: 000040 曾在订单头新增秒杀列, 015 评审修复轮以 000041 将 `bargain_record.order_no` 改可空（C1）, 000042 撤销订单头死列、改用 `trade_order_item.flash_sale_item_id`（000023 既有, 行级归属）。

## 一、五类活动的存储（本批只读；秒杀另有写入）

| 表 | 本批用法 | 关键列 |
|---|---|---|
| `flash_sale_activity` | 列表（进行中+预告） | 起止时间、status、deleted |
| `flash_sale_item` | 列表 + **下单取价与扣活动库存** | `flash_price`、`stock_count`、`sold_count`、`per_limit`、**`uk_activity_sku(activity_id,sku_id)`** |
| `group_buy_activity` / `group_buy_item` | 列表 | 成团人数、`group_price` |
| `bargain_activity` / `bargain_item` | 列表 + 发起 | `original_price`、`floor_price`、`max_cut_count`、`config` |
| `assist_activity` | 列表 + 发起 | `required_count`、`per_limit`、`reward_type`(1券/2积分)、`reward_ref` |
| `promotion_activity` / `_ladder` / `_scope` | 满减列表 | 档位与范围（结算侧批次 06 已消费） |

## 二、玩法的写入与状态机

### 砍价（`bargain_record` + `bargain_helper`，helper 只追加）

```
发起（POST /bargains）: 校验场次商品有效 + 活动时间内 + 发起次数未超 → 建单
  current_price = original_price − 每刀金额（D4 均分，向下取整到分）
  cut_count = 1（发起即首刀）; status = 1 砍价中; expire_time = now + 活动配置/默认时长
帮砍（POST /bargains/{id}/cut）:
  校验: 单存在/未超时/发起者本人不可帮砍/该人未帮砍过（唯一键兜底）
  条件更新: 当前价 − 每刀金额（不足则补齐到 floor_price）; cut_count +1
  到底价 → status = 2 到底价（可下单）
超时: 查询时惰性判定（expire_time < now 且 status=1 → 视作 4 超时）或由后续批次调度置位
```

| status | 含义 | 可帮砍 | 可下单 |
|---|---|---|---|
| 1 | 砍价中 | ✅ | ❌（未到底价） |
| 2 | 到底价 | ❌（已到底） | ✅ |
| 3 | 已下单 | ❌ | ❌ |
| 4 | 超时 | ❌ | ❌ |
| 5 | 取消 | ❌ | ❌ |

### 助力（`assist_record` + `assist_helper`，helper 只追加）

```
发起（POST /assists）: 校验活动在时间内 + 该人发起次数 < per_limit → 建记录（helper_count=0, status=1 进行中）
助力（POST /assists/{id}/helpers）: 校验 非本人/未助力过（唯一键）/活动未结束 → 写 helper + helper_count+1（条件更新）
  helper_count >= required_count → status = 2 已完成发奖 + finish_time + **端口投递发奖意图**（D3）
超时: 活动结束而未达标 → 查询时视作 3 已过期
```

| status | 含义 |
|---|---|
| 1 | 进行中 |
| 2 | 已完成发奖 |
| 3 | 已过期 |

## 三、秒杀下单的双库存（跨 `flash_sale_item` 与 `inventory`）

| 步骤 | 条件更新 | 失败 |
|---|---|---|
| 活动库存 | `WHERE id=? AND status=1 AND NOW() BETWEEN 起止 AND stock_count - sold_count >= 数量` → `sold_count += 数量` | 40003 已抢完 |
| 商品库存 | `WHERE sku_id=? AND total - locked >= 数量` → `locked += 数量` | 40001 库存不足 |
| 取消回补 | `sold_count -= 数量` + `locked -= 数量`（各自条件防负） | 记日志（不阻断取消） |
| 支付核销 | `total -= 数量`、`locked -= 数量`（**既有逻辑，本批不改**） | 回滚（批次 07 已实现） |

**场次商品校验**: 按 `(activity_id, sku_id)` 二元组查 `flash_sale_item`（唯一键）→ 命中才允许使用该秒杀价；未命中即拒绝（防篡改）。

## 四、新增 DTO（本批新建，管理面 DTO 不动）

| DTO | 用途 |
|---|---|
| `PublicGroupBuyItem` / `PublicFlashSaleItem` / `PublicBargainItem` / `PublicAssistItem` / `PublicFullReductionItem` | 五类公开列表项（含活动信息 + 场次商品摘要） |
| `ActivitySkuBrief`（或复用既有） | 场次商品摘要（skuId/价/剩余/限购） |
| `BargainLaunchResult` / `BargainProgressView` / `PlayHelperItem` | 砍价：发起出参 / 进度视图 / 帮砍条目（昵称脱敏） |
| `AssistLaunchResult` / `AssistProgressView` | 助力：发起出参 / 进度视图（含助力条目） |
| `IndexAggregate` | 首页聚合（轮播 + 楼层 + 五入口 + 券） |

## 五、校验规则汇总（可测）

| 规则 | 出处 | 失败码 |
|---|---|---|
| 秒杀：场次商品 SKU 匹配 | FR-002 | 50002 活动无效 |
| 秒杀：活动库存充足 | FR-003 | 40003 已抢完 |
| 秒杀：商品库存充足 | FR-003 | 40001 库存不足 |
| 秒杀：场次在时间窗内/已启用 | FR-004 | 50002 |
| 砍价：本人不可帮砍 / 一人一刀 | FR-009 | 50004 已帮砍过 / 50002 |
| 砍价：不砍穿（≥ 底价） | FR-009 | —（算法保证） |
| 砍价：超时不可砍 | FR-010 | 50002 活动无效 |
| 砍价：发起次数上限 | FR-007 | 50006 发起次数用尽 |
| 助力：本人不可助力 / 一人一助力 | FR-013 | 50007 已助力 / 50002 |
| 助力：发起次数上限 | FR-012 | 50006 |
| 助力：活动时间窗 | FR-012 | 50002 |
| 风控拦截 | FR-015 | 10005 无权限（沿用风控既有语义） |

**错误码（既有，无需新增）**: `40001 库存不足`、`40003 已抢完`、`50002 活动无效`、`50004 已帮砍过`、`50006 发起次数用尽`、`50007 已助力`。
