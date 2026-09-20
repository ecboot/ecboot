# shop 渠道契约（商城 C 端）

前缀 `/shop`。游客浏览 + 会员交易混合，逐端点标注鉴权。

## 商品浏览（公开）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/shop/categories` | 公开 | 三级分类树 | — | tree[id,name,icon,children] | — |
| GET | `/shop/brands` | 公开 | 品牌列表（启用） | page? | list[id,name,logo] | — |
| GET | `/shop/products` | 公开 | 商品列表（在售） | categoryId?, brandId?, sort(0综合/1销量/2价格/3上新), priceMin?,priceMax?, page | list[spuId,spuName,image,priceRange,saleCount], total | — |
| GET | `/shop/products/{spuId}` | 公开 | 商品详情 | — | spu 全字段(无成本), skus[id,skuNo,specs,price,linePrice,sellable], specDefinitions, freight 概要, reviewSummary{avg,distribution} | 30001 商品不存在/已下架 |
| GET | `/shop/search` | 公开 | 关键词搜索 | keyword, categoryId?, brandId?, sort, page | 同商品列表 | — |
| GET | `/shop/products/{spuId}/reviews` | 公开 | 商品评价列表（审核通过） | score?, page | list[reviewId,user(匿名处理),score,content,images,specs,reply,extra,createdAt], total, summary | — |

## 营销活动（公开浏览）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/shop/activities/group-buys` | 公开 | 进行中拼团列表 | page | list[activityId,name,spu 摘要,items[{skuId,groupPrice}],groupSize,endTime] | — |
| GET | `/shop/activities/flash-sales` | 公开 | 进行中/预告秒杀列表 | page | list[activityId,name,timeRange,items[{skuId,flashPrice,stockRemain,perLimit}]] | — |
| GET | `/shop/activities/bargains` | 公开 | 进行中砍价列表 | page | list[activityId,name,spu 摘要,items[{skuId,originalPrice,floorPrice}]] | — |
| GET | `/shop/activities/assists` | 公开 | 进行中助力列表 | page | list[activityId,name,rewardDesc,requiredCount,endTime] | — |
| GET | `/shop/full-reductions` | 公开 | 当前满减活动（含档位） | spuId?（范围命中过滤） | list[activityId,name,ladders[{threshold,discount}]] | — |

## 运营位（公开）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/shop/banners` | 公开 | 轮播/弹窗（投放时段内） | position | list[id,imageUrl,linkUrl] | — |
| GET | `/shop/floors` | 公开 | 楼层内容（含商品摘要装配） | — | list[floorId,floorType,title,items[]] | — |

## 购物车（会员）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/shop/cart` | 会员 | 购物车（含实时价态与可售态） | — | list[itemId,skuId,spuName,specs,price,sellable,quantity,checked] | 10003 |
| POST | `/shop/cart/items` | 会员 | 加购（重复=数量累加） | skuId, quantity | itemId | 30002 SKU 不可售 |
| PUT | `/shop/cart/items/{itemId}` | 会员 | 改数量/勾选 | quantity?, checked? | true | — |
| DELETE | `/shop/cart/items/{itemId}` | 会员 | 移除 | — | true | — |
| GET | `/shop/cart/checkout` | 会员 | 结算试算（勾选项） | addressId?, couponId?, usePoint?, useAccount? | 商品金额/优惠明细(券/满减/积分/余额)/运费/应付, 可用券列表 | 30002, 30003 部分商品失效 |

## 订单（会员）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/shop/orders` | 会员 | 创建订单（幂等） | addressId, requestToken, userCouponId?, usePoint?, useAccount?, userRemark?；玩法上下文二选一：cartItems(普通) / groupBuyTeamId / flashSaleItemId / bargainRecordId | orderNo, payAmount | 40001 库存不足, 40002 券不可用, 40003 已抢完, 40004 砍价单不可下单, 10004 幂等冲突 |
| GET | `/shop/orders` | 会员 | 订单列表 | status?, page | list[orderNo,status,amount 摘要,items 摘要,createdAt], total | 10003 |
| GET | `/shop/orders/{orderNo}` | 会员 | 订单详情（含状态流转时间线） | — | 全字段+items+pay 摘要+deliver 信息+cancel 信息 | 40005 订单不存在 |
| POST | `/shop/orders/{orderNo}/cancel` | 会员 | 取消（仅待付款） | reason? | true | 40006 状态不允许 |
| POST | `/shop/orders/{orderNo}/confirm` | 会员 | 确认收货（仅待收货） | — | true | 40006 |

## 支付（会员 + 开放回调）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/shop/pay` | 会员 | 发起支付（创建支付单，返回渠道唤起参数） | orderNo, payChannel? | payNo, channelParams | 40007 支付单创建失败 |
| GET | `/shop/pay/status` | 会员 | 支付状态查询 | payNo | status, paidAt | — |
| POST | `/shop/pay/notify` | 开放（签名验证） | 支付结果回调（渠道→平台） | 渠道原文 | 渠道应答格式 | — |
| POST | `/shop/refund/notify` | 开放（签名验证） | 退款结果回调 | 渠道原文 | 渠道应答格式 | — |

## 售后（会员）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/shop/after-sales` | 会员 | 申请售后（按订单项） | orderItemId, type(1仅退款/2退货退款), quantity, reason, description?, voucherImages? | afterSaleNo | 40008 不可售后（状态/超额） |
| GET | `/shop/after-sales` | 会员 | 售后列表 | status?, page | list[afterSaleNo,type,refundAmount,status,createdAt], total | 10003 |
| GET | `/shop/after-sales/{afterSaleNo}` | 会员 | 售后详情（含进度时间线） | — | 全字段 | 40009 |
| POST | `/shop/after-sales/{afterSaleNo}/cancel` | 会员 | 撤销申请 | — | true | 40006 |
| POST | `/shop/after-sales/{afterSaleNo}/logistics` | 会员 | 填写寄回单号（type=2） | returnLogisticsNo | true | 40006 |

## 评价（会员 + 公开列表在商品接口）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/shop/reviews` | 会员 | 提交评价（一项一评） | orderItemId, score(1-5), content?, images?, isAnonymous? | reviewId | 40010 已评价 |
| POST | `/shop/reviews/{reviewId}/extra` | 会员 | 追评（一次，90 天内） | content, images? | true | 40011 已追评/超期 |
| GET | `/shop/reviews/mine` | 会员 | 我的评价 | page | list[reviewId,spuName,score,content,status,createdAt], total | 10003 |

## 砍价 / 助力（会员发起，公开查看进度；帮砍/助力=风控挂载位）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/shop/bargains` | 会员 | 发起砍价 | bargainItemId | bargainRecordId, currentPrice | 50002 活动无效 |
| GET | `/shop/bargains/{recordId}` | 公开 | 砍价进度（含帮砍列表） | — | currentPrice,cutCount,status,helpers[],expireTime | 50003 |
| POST | `/shop/bargains/{recordId}/cut` | 会员 | 帮砍（一人一刀；**风控检查位**） | — | cutAmount, currentPrice | 50004 已砍过, 10004 频控, 50005 已到底价 |
| POST | `/shop/assists` | 会员 | 发起助力 | activityId | assistRecordId | 50006 次数用尽 |
| GET | `/shop/assists/{recordId}` | 公开 | 助力进度 | — | helperCount,requiredCount,status | 50003 |
| POST | `/shop/assists/{recordId}/helpers` | 会员 | 助力（一人一次；**风控检查位**） | — | true | 50007 已助力 |
