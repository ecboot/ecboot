# Data Model: 营销后台（016-marketing-admin）

**全部为既有表**（000005/000017/000018/000020/000029），本批零迁移。只列本批消费的字段与约束。

## coupon / user_coupon（券, 000005 + 批次 06 消费端）

- `coupon`: name/type(1满减 2无门槛)/threshold/discount/total_count(0不限)/received_count/per_limit/valid_type(1固定 2领取后N天)/valid_start_at/valid_end_at/valid_days/status(1启用 0停发)/deleted
- `user_coupon`: coupon_id/user_id/status(1未用 2已用 3过期 4退回)/order_no/created_at —— AdminRecords 读侧
- 校验: type=1 须 threshold；validType 两型互斥条件；discount/threshold 金额 ≤2 位小数（money.FromYuanString）

## promotion_activity / _ladder / _scope（满减, 000020）

- `_ladder`: uk_activity_threshold(activity_id, threshold_amount) → 50008；列表档位按 threshold 升序
- `_scope`: scope_type(1全场 2分类 3商品)/target_id(全场 NULL)；"空 scopes=全场"——创建不落行, 详情回读空数组口径="全场"（与 C 端 scopeDescOf 一致）
- 更新: 档位+范围**全量替换**（事务先删后插）

## group_buy_activity / _item（拼团, 000017/000024）

- 活动: spu_id/group_size(≥2)/per_limit/时间窗/status/deleted
- `_item`: uk_activity_sku(activity_id, sku_id)/group_price —— SetItems 全量替换

## flash_sale_activity / _item（秒杀, 000018）

- `_item`: uk_activity_sku/flash_price/stock_count/sold_count/per_limit + CHECK(sold_count≤stock_count)
- **已售保护**: 行被 `trade_order_item.flash_sale_item_id` 引用且 sold_count>0 → 禁止移除（50002）；替换事务内锁定读复核

## bargain_activity / _item（砍价, 000029）

- 活动: spu_id/时间窗; `_item`: original_price(起始价)/floor_price(底价)/max_cut_count(0=不限)/config
- 校验: original>floor>0；max_cut_count≥0

## assist_activity（助力, 000029）

- reward_type(1券 2积分)/reward_ref(券ID; 积分=0)/required_count(≥1)/per_limit(>1 不生效, uk_activity_user 钉死每人一次)/config(JSON: pointAmount)/status/deleted
