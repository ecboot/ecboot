# ecboot-service-shop（商城领域服务）

**商品与交易域**的领域实现——从商品目录到资金闭环的完整交易链。下一步按领域拆分、内部 **DDD 分层**实现（模块边界与对外契约不变）。

## 包结构约定（DDD 四层）

```
org.juling.ecboot.shop
├── interfaces/        # 对外暴露：供渠道调用的应用接口 + 领域事件（订单已支付/售后完成等）
├── application/       # 用例层：下单编排（锁库存+优惠+幂等）、支付回调、售后流转（事务边界）
├── domain/            # 领域层：订单聚合（状态机）、库存（三段式）、促销策略、评价聚合、仓储接口
└── infrastructure/    # 基础设施：JPA 仓储实现、支付渠道适配、Redis 库存挡板、ES 同步
```

## 功能内容（对齐 Schema 表）

| 功能域 | 内容 | 对应表 |
|---|---|---|
| 商品目录 | 三级分类树、品牌、SPU（规格定义 JSON/图集/上下架）、SKU（规格快照/两级上下架） | `product_category/brand/spu/sku` |
| 库存 | **三段式锁定模型**（下单锁/支付核销/取消释放/售后回补，条件更新防超卖，ADR-0001）、全量流水与对账 | `inventory`、`inventory_log` |
| 运费 | 运费模板（按件/按重、省份差异化、满额包邮）、下单快照 | `freight_template`、`freight_rule` |
| 购物车 | 加购/改量/勾选、结算校验（实时价格与可售） | `cart_item` |
| 交易（核心） | 订单状态机（10/20/30/40/90）、**订单项+收货+运费快照**、下单幂等（request_token）、优惠三构成恒等式、超时取消（30min）/自动收货（7 天） | `trade_order`、`trade_order_item`、`trade_order_log` |
| 支付 | 支付单（一单多尝试）、回调幂等（条件更新+渠道单号唯一+金额校验）、原文留档 | `pay_order`、`pay_callback_log` |
| 售后 | 仅退款/退货退款状态机、按订单项粒度、退款渠道幂等（out_refund_no）、完成回补库存+联动佣金冲销（领域事件） | `after_sale_order` |
| 促销 | 优惠券模板（防超发）、满减（多档位/范围关系表/先满减后券）、拼团（一团一单/超时解散）、秒杀（**活动库存分账**、条件更新不超卖） | `coupon`、`promotion_activity/ladder/scope`、`group_buy_*`、`flash_sale_*` |
| 评价 | 一项一评、商品快照、追评/商家回复各一次、审核 | `product_review` |
| 物流 | 物流公司字典维护（编码唯一/停用保留） | `logistics_company` |
| 风控 | 规则与事件（黑名单/高频/套利特征）、申诉流转 | `risk_rule`、`risk_record` |

## 职责边界

- **做**：交易闭环全部业务规则；发布领域事件（订单已支付→user 域计佣；售后完成→佣金冲销）
- **不做**：不暴露 REST；不依赖 api 渠道/兄弟服务/start（enforcer）；会员/积分/关系链在 user 域（券实例由 user 域持有，本域只管模板与核销协议）

## 依赖关系

- 白名单：`ecboot-common`、`ecboot-infra-core`（兄弟 service-user 互禁，经领域事件协作）
- 被依赖：api-shop / api-admin 渠道

相关文档：ADR-0001（库存模型）、`docs/schema-design.md`（状态机/幂等总则/恒等式）
