# Data Model: 全角色 API 接口面（Phase 1）

本特性**不新增业务实体与表**——"数据模型"即契约层结构模型：通用请求/响应结构、鉴权模型、权限点清单、错误码登记表。业务数据模型 = 既有 68 表（见 `docs/schema-design.md`）。

## 一、通用契约结构（api 层共享定义）

| 结构 | 字段 | 说明 |
|---|---|---|
| 统一响应包装 | code(int) / message(string) / data(T) | 中间件包装，Res 结构体仅定义 data |
| 分页请求 | page(默认1) / pageSize(默认10,≤100) | 嵌入各列表 Req |
| 分页响应 | total(int64→string) / list([]T) | 同上 |
| 身份上下文 | userId / adminId | 中间件注入（ctx 取用），客户端不可传 |
| 金额 | string（元，两位小数） | 禁 numeric JSON（精度） |
| ID | string（int64 字符串化） | 同上 |
| 时间 | RFC3339 字符串 | — |

## 二、鉴权模型

| 组 | 路由前缀 | 中间件语义 |
|---|---|---|
| 公开组 | `/common/*`、`/shop` 浏览类、`/user/login/*`、`/admin/login`、回调 `/shop/pay/notify`、`/shop/refund/notify` | 无 / 回调签名验证 |
| 会员组 | `/user/*`（除登录）、`/shop` 交易类 | Bearer token 校验 + 滑动续期（003 会话） |
| 管理组 | `/admin/*`（除登录） | Bearer 后台凭证 + 权限点校验（is_super 直通） |

## 三、权限点清单（RBAC 种子，`admin_permission.code`）

按模块前缀登记（`模块:资源:动作`，动作 ⊆ read/create/update/delete + 专属动作）：

- `product:{category,brand,spu,sku}:{read,create,update,delete}`（spu:update 含上下架/限售）
- `inventory:{read,adjust}`
- `order:{read,deliver,cancel,update}`
- `aftersale:{read,audit,refund}`
- `promotion:{coupon,fullreduction,groupbuy,flashsale,bargain,assist}:{read,create,update,delete}`
- `operation:{banner,floor}:{read,create,update,delete}`；`store:manage:{...}`；`logistics:company:{...}`
- `distribution:{read,audit,rule:*,withdraw:read,withdraw:audit,withdraw:pay}`
- `member:{read,update}`
- `risk:{rule:*,record:read,record:appeal}`
- `system:{role:*,admin:*,audit:read,config:read,config:update}`
- `dashboard:read`

## 四、错误码登记表（本特性新增段位；具体码值随实现登记 errcode 包）

| 段 | 域 | 已在契约中出现的语义 |
|---|---|---|
| 3xxxx | 商品 | 30001 不存在/下架、30002 SKU 不可售、30003 商品失效、30004 分类有引用、30005 无启用 SKU、30006 规格组合重复、30007 库存调整非法 |
| 4xxxx | 交易/售后/支付 | 40001 库存不足、40002 券不可用、40003 已抢完、40004 砍价单不可下单、40005 订单不存在、40006 状态不允许、40007 支付单失败、40008 不可售后、40009 售后单不存在、40010 已评价、40011 追评违规 |
| 5xxxx | 促销 | 50001 券领完/超限、50002 活动无效、50003 活动/单不存在、50004 已砍过、50005 已到底价、50006 次数用尽、50007 已助力、50008 档位门槛重复 |
| 6xxxx | 分销 | 60001 已申请、60002 余额不足、60003 进行中存在单 |
| 7xxxx | 门店 | 70001 不存在/歇业 |
| 8xxxx | 后台治理 | 80001 凭证错误、80002 原密码错误 |

## 五、接口定义落位（api 目录文件 → 契约域）

| 渠道 | v1 文件 | 覆盖契约 |
|---|---|---|
| common | captcha.go / store.go / ping.go / share.go | common-api 全部 |
| user | auth.go（含 refresh）/ profile.go / address.go / favorite.go / footprint.go / point.go / coupon.go / message.go / distribution.go | user-api 全部 |
| shop | product.go / activity.go / banner.go / cart.go / order.go / pay.go / aftersale.go / review.go / bargain.go / assist.go | shop-api 全部 |
| admin | auth.go / product.go / inventory.go / order.go / aftersale.go / promotion.go / operation.go / store.go / logistics.go / distribution.go / member.go / rbac.go / audit.go / risk.go / config.go / dashboard.go | admin-api 全部 |
