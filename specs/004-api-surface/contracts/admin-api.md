# admin 渠道契约（管理后台）

前缀 `/admin`。全部管理员鉴权；写操作标注权限点（`模块:资源:动作`，RBAC 种子来源）。admin 登录属本特性（003 仅覆盖 C 端）。

## 认证与账号（双凭证同会员侧）

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| POST | `/admin/login` | 公开 | — | 后台登录（含登录审计） | username, password(或+图形码) | token, refreshToken, admin{realName,isSuper} | 20003, 80001 凭证错误 |
| POST | `/admin/token/refresh` | 公开（凭 refreshToken） | — | 刷新凭证 | refreshToken | token, refreshToken | 10003 |
| POST | `/admin/logout` | 管理员 | — | 登出 | — | true | — |
| GET | `/admin/profile` | 管理员 | — | 个人信息 | — | username, realName, roles[] | — |
| PUT | `/admin/profile/password` | 管理员 | — | 改密 | oldPassword, newPassword | true | 80002 原密码错误 |

## 商品管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/categories` | 管理员 | product:category:read | 分类树（含禁用） | — | tree | — |
| POST | `/admin/categories` | 管理员 | product:category:create | 新增分类 | parentId,name,level,icon?,sort | id | — |
| PUT/DELETE | `/admin/categories/{id}` | 管理员 | product:category:update/delete | 改/删（删=软删，有商品禁删） | — | true | 30004 有商品引用 |
| GET/POST | `/admin/brands` | 管理员 | product:brand:read/create | 品牌列表/新增 | — | — | — |
| PUT/DELETE | `/admin/brands/{id}` | 管理员 | product:brand:update/delete | 改/删 | — | true | — |
| GET | `/admin/products` | 管理员 | product:spu:read | SPU 列表（全状态筛选） | status?,categoryId?,keyword?,page | list+total | — |
| POST | `/admin/products` | 管理员 | product:spu:create | 创建 SPU（含规格定义/图文/运费模板） | name,categoryId,brandId?,images,specDefinitions?,attributes?,freightTemplateId? | spuId, spuNo | — |
| GET/PUT/DELETE | `/admin/products/{spuId}` | 管理员 | product:spu:read/update/delete | 详情/改/删（软删） | — | — | 30001 |
| POST | `/admin/products/{spuId}/status` | 管理员 | product:spu:update | 上下架 | status(0/1) | true | 30005 无启用 SKU 禁上架 |
| PUT | `/admin/products/{spuId}/restrict-codes` | 管理员 | product:spu:update | 限售区域设置 | saleRestrictCodes[]或 null | true | — |
| POST | `/admin/products/{spuId}/skus` | 管理员 | product:sku:create | 新增 SKU（规格组合冲突拒绝） | specs,price,linePrice?,image?,weight?,barcode? | skuId, skuNo | 30006 规格组合重复 |
| PUT/DELETE | `/admin/skus/{skuId}` | 管理员 | product:sku:update/delete | 改/删 SKU | — | true | — |
| POST | `/admin/skus/{skuId}/status` | 管理员 | product:sku:update | SKU 启停 | status | true | — |

## 库存管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/inventories` | 管理员 | inventory:read | 库存列表（可售/锁定） | skuId?,keyword?,page | list[skuId,skuName,total,locked,available], total | — |
| POST | `/admin/inventories/{skuId}/adjust` | 管理员 | inventory:adjust | 库存调整（留痕 inventory_log） | delta(±), remark | totalAfter | 30007 调整后非法 |
| GET | `/admin/inventories/warnings` | 管理员 | inventory:read | 预警列表（≤warn_count） | page | list+total | — |

## 订单管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/orders` | 管理员 | order:read | 订单列表（状态/时间/用户/单号筛选） | status?,userKeyword?,orderNo?,timeRange?,page | list+total | — |
| GET | `/admin/orders/{orderNo}` | 管理员 | order:read | 订单详情（含收货快照/状态流水） | — | 全字段+items+logs | 40005 |
| POST | `/admin/orders/{orderNo}/deliver` | 管理员 | order:deliver | 发货 | logisticsCode, deliverNo | true | 40006 状态不允许 |
| POST | `/admin/orders/{orderNo}/cancel` | 管理员 | order:cancel | 管理员取消（待付款） | reason | true | 40006 |
| PUT | `/admin/orders/{orderNo}/seller-remark` | 管理员 | order:update | 卖家备注 | sellerRemark | true | — |

## 售后管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/after-sales` | 管理员 | aftersale:read | 售后列表（状态筛选） | status?,page | list+total | — |
| GET | `/admin/after-sales/{afterSaleNo}` | 管理员 | aftersale:read | 售后详情 | — | 全字段 | 40009 |
| POST | `/admin/after-sales/{afterSaleNo}/approve` | 管理员 | aftersale:audit | 同意（仅退款→待退款；退货→待寄回） | — | true | 40006 |
| POST | `/admin/after-sales/{afterSaleNo}/reject` | 管理员 | aftersale:audit | 拒绝 | rejectReason | true | 40006 |
| POST | `/admin/after-sales/{afterSaleNo}/confirm-receipt` | 管理员 | aftersale:audit | 退货确认收货（→待退款） | — | true | 40006 |
| POST | `/admin/after-sales/{afterSaleNo}/retry-refund` | 管理员 | aftersale:refund | 退款重试（失败后） | — | true | 40006 |

## 促销管理（优惠券/满减/拼团/秒杀/砍价/助力）

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET/POST | `/admin/coupons` | 管理员 | promotion:coupon:* | 券模板列表/创建 | 创建：name,type,threshold,discount,totalCount,perLimit,valid 规则 | list+total / couponId | — |
| GET/PUT/DELETE | `/admin/coupons/{id}` | 管理员 | promotion:coupon:* | 详情/改/停发 | — | — | — |
| GET | `/admin/coupons/{id}/records` | 管理员 | promotion:coupon:read | 发放/使用记录 | page | list+total | — |
| GET/POST | `/admin/full-reductions` | 管理员 | promotion:fullreduction:* | 满减活动列表/创建（档位+范围嵌套提交） | name,timeRange,ladders[],scopes[] | activityId | 50008 档位门槛重复 |
| GET/PUT/DELETE | `/admin/full-reductions/{id}` | 管理员 | promotion:fullreduction:* | 详情/改/停用 | — | — | — |
| GET/POST | `/admin/group-buys` … `/admin/group-buys/{id}` 及 `/{id}/items` | 管理员 | promotion:groupbuy:* | 拼团活动 CRUD + SKU 场次商品管理 | name,spuId,timeRange,groupSize,perLimit; items[{skuId,groupPrice}] | — | — |
| GET/POST | `/admin/flash-sales` … 同构 | 管理员 | promotion:flashsale:* | 秒杀活动 CRUD + 场次商品（限量/秒杀价/限购） | — | — | — |
| GET/POST | `/admin/bargains` … 同构 | 管理员 | promotion:bargain:* | 砍价活动 CRUD + SKU 价格区间 | — | — | — |
| GET/POST | `/admin/assists` … `/admin/assists/{id}` | 管理员 | promotion:assist:* | 助力活动 CRUD | rewardType,rewardRef,requiredCount,perLimit,timeRange | — | — |

## 运营位 / 门店 / 物流

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET/POST/PUT/DELETE | `/admin/banners`、`/admin/banners/{id}` | 管理员 | operation:banner:* | 轮播管理 | position,imageUrl,linkUrl,sort,timeRange,status | — | — |
| GET/POST/PUT/DELETE | `/admin/floors`、`/admin/floors/{id}` | 管理员 | operation:floor:* | 楼层管理 | floorType,title,config,sort,status | — | — |
| GET/POST/PUT/DELETE | `/admin/stores`、`/admin/stores/{id}` | 管理员 | store:manage:* | 门店管理（含经纬度/自提/营业状态） | name,区划码,detail,longitude,latitude,businessHours,pickupEnabled,status | — | — |
| GET/POST/PUT/DELETE | `/admin/logistics-companies`、`/{id}` | 管理员 | logistics:company:* | 物流公司字典 | code,name,trackingRule,status | — | — |

## 分销管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/distributors` | 管理员 | distribution:read | 推广员列表（状态筛选） | status?,page | list+total | — |
| POST | `/admin/distributors/{id}/audit` | 管理员 | distribution:audit | 审核（通过/拒绝） | pass(bool) | true | — |
| POST | `/admin/distributors/{id}/freeze` | 管理员 | distribution:audit | 冻结/解冻 | freeze(bool) | true | — |
| GET/POST | `/admin/commission-rules` | 管理员 | distribution:rule:* | 佣金规则列表/创建 | scopeType,scopeId,level1Rate,level2Rate | ruleId | — |
| PUT/DELETE | `/admin/commission-rules/{id}` | 管理员 | distribution:rule:* | 改/删 | — | true | — |
| GET | `/admin/commission-records` | 管理员 | distribution:read | 佣金记录（全局） | status?,orderNo?,page | list+total | — |
| GET | `/admin/withdraws` | 管理员 | distribution:withdraw:read | 提现列表 | status?,page | list+total | — |
| POST | `/admin/withdraws/{withdrawNo}/audit` | 管理员 | distribution:withdraw:audit | 提现审核（通过/拒绝；拒绝回退余额） | pass(bool), reason? | true | 40006 |
| POST | `/admin/withdraws/{withdrawNo}/pay` | 管理员 | distribution:withdraw:pay | 标记打款/打款结果登记 | success(bool), channelOrderNo?, failReason? | true | 40006 |
| GET | `/admin/invite-records` | 管理员 | distribution:read | 邀请激励记录 | page | list+total | — |

## 会员管理

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET | `/admin/members` | 管理员 | member:read | 会员列表（手机号**精确检索**、脱敏展示） | phone?,status?,page | list+total | — |
| GET | `/admin/members/{userId}` | 管理员 | member:read | 会员详情（含资产概要） | — | 全字段（脱敏） | — |
| POST | `/admin/members/{userId}/disable` | 管理员 | member:update | 禁用/启用 | disable(bool), reason? | true | — |
| POST | `/admin/members/{userId}/rebind-phone` | 管理员 | member:update | 改绑手机号（旧号解占，审计留痕） | newPhone | true | — |

## 治理（RBAC / 审计 / 风控 / 配置 / 看板）

| 方法 | 路径 | 鉴权 | 权限点 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|---|
| GET/POST | `/admin/roles`、GET/PUT/DELETE `/admin/roles/{id}` | 管理员 | system:role:* | 角色 CRUD | name,code,description | — | — |
| GET | `/admin/permissions` | 管理员 | system:role:read | 权限树 | — | tree | — |
| PUT | `/admin/roles/{id}/permissions` | 管理员 | system:role:assign | 角色-权限全量替换 | permissionIds[] | true | — |
| GET/POST | `/admin/admin-users`、GET/PUT/DELETE `/admin/admin-users/{id}` | 管理员 | system:admin:* | 后台账号 CRUD | username,password(创建),realName,status | — | — |
| PUT | `/admin/admin-users/{id}/roles` | 管理员 | system:admin:assign | 账号-角色分配 | roleIds[] | true | — |
| GET | `/admin/operation-logs` | 管理员 | system:audit:read | 操作审计 | adminId?,module?,timeRange?,page | list+total | — |
| GET | `/admin/admin-login-logs` | 管理员 | system:audit:read | 后台登录审计 | username?,page | list+total | — |
| GET/POST | `/admin/risk-rules`、GET/PUT/DELETE `/{id}` | 管理员 | risk:rule:* | 风控规则 CRUD | name,ruleType,conditionExpr,action | — | — |
| GET | `/admin/risk-records` | 管理员 | risk:record:read | 风控事件 | userId?,appealStatus?,page | list+total | — |
| POST | `/admin/risk-records/{id}/appeal` | 管理员 | risk:record:appeal | 申诉处理 | pass(bool), remark | true | — |
| GET | `/admin/configs` | 管理员 | system:config:read | 配置列表 | — | list[code,value,valueType,name,description,status] | — |
| PUT | `/admin/configs/{code}` | 管理员 | system:config:update | 修改配置 | value, status? | true | — |
| GET | `/admin/dashboard/trade` | 管理员 | dashboard:read | 交易看板 | timeRange | 订单数/销售额/退款额 | — |
| GET | `/admin/dashboard/member` | 管理员 | dashboard:read | 会员看板 | timeRange | 新增/活跃/休眠数 | — |
| GET | `/admin/dashboard/product` | 管理员 | dashboard:read | 商品看板 | — | 在售数/低库存预警数 | — |
