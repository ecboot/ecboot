# ecboot-api-shop（商城渠道）

小程序前台·**商城** REST API——逛、买、付、售后的完整 C 端交易链端点。

## 功能内容（REST 端点规划）

| 功能组 | 端点（规划） | 后端域 |
|---|---|---|
| 商品浏览 | 分类树、品牌列表、SPU 列表（分类/筛选/排序）、SPU 详情（SKU+规格+评价汇总） | service-shop · 目录 |
| 搜索 | 关键词搜索（ES，名称/品牌/价格区间/销量与价格排序） | service-shop |
| 购物车 | 加购/改量/移除、勾选结算、结算页（实时价态/优惠试算/运费试算） | service-shop |
| 下单交易 | 创建订单（幂等 token；普通/拼团/秒杀）、订单列表/详情、取消、确认收货 | service-shop · 交易 |
| 支付 | 发起支付（多渠道尝试）、支付结果查询 | service-shop · 支付 |
| 售后 | 申请（仅退款/退货退款）、撤销、填写寄回单号、进度查询 | service-shop · 售后 |
| 评价 | 提交评价（一项一评）、追评、商品评价列表 | service-shop · 评价 |
| 优惠券 | 可领模板列表、下单可用券匹配 | service-shop/user（模板/持有协作） |
| 秒杀/拼团 | 活动列表与详情、参团/开团、我的团 | service-shop · 促销 |

## 职责边界

- **做**：C 端交易 REST 暴露（薄控制器），编排 service-shop；支付回调端点亦在此渠道
- **不做**：与 api-common/user/admin **互不编译依赖**；会员资产端点在 api-user；后台管理在 api-admin

## 依赖关系

- 白名单：`ecboot-api-webmvc`、`ecboot-service-user`、`ecboot-service-shop`、`ecboot-common`、`ecboot-infra-core`
- 被依赖：仅 `ecboot-start`

相关文档：`../README.md`、`ecboot-service-shop/README.md`（状态机/幂等总则）
