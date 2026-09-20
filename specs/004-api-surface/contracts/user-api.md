# user 渠道契约（会员中心）

前缀 `/user`。鉴权=会员（除 auth 登录与公开标注外）；双凭证会话。

## 认证与会话（003 设计的路径迁移版）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/user/login/sms` | 公开 | 短信验证码登录（注册即登录） | phoneNumber, smsCode, channel? | token, refreshToken, userId, isNew | 10004, 20001, 20002, 20003 |
| POST | `/user/login/wx` | 公开 | 微信登录（手机号优先归并；mock） | wxCode, phone?, smsCode? | 同上 | 20001, 20003, 20004 |
| POST | `/user/token/refresh` | 公开（凭 refreshToken） | 刷新访问凭证 | refreshToken | token, refreshToken | 10003 |
| POST | `/user/logout` | 会员 | 登出（双凭证同失效） | — | true | 10003 |
| GET | `/user/me` | 会员 | 当前用户信息 | — | userId, nickname, avatar, phone(脱敏), level, registerChannel | 10003 |
| GET | `/user/login-logs` | 会员 | 近 30 天登录记录 | page | list[time,channel,status,ip(脱敏)], total | 10003 |

## 个人资料 / 等级

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/user/profile` | 会员 | 资料详情 | — | nickname, avatar, gender, growthValue, level{name,benefits} | 10003 |
| PUT | `/user/profile` | 会员 | 修改资料 | nickname?, avatar?, gender? | true | 10001 |
| GET | `/user/point-account` | 会员 | 积分账户 | — | balance, lastEarnedAt | 10003 |
| GET | `/user/point-logs` | 会员 | 积分流水 | bizType?, page | list[points,balanceAfter,bizType,orderNo,createdAt], total | 10003 |

## 收货地址

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/user/addresses` | 会员 | 地址列表 | — | list[id,receiverName,receiverPhone(脱敏),province..detail,provinceCode..districtCode,isDefault] | 10003 |
| POST | `/user/addresses` | 会员 | 新增地址 | receiverName, receiverPhone, 省市区+区划码, detailAddress, isDefault? | id | 10001 |
| PUT | `/user/addresses/{id}` | 会员 | 修改地址 | 同上（部分可选） | true | 10001 |
| DELETE | `/user/addresses/{id}` | 会员 | 删除地址（软删） | — | true | — |
| PUT | `/user/addresses/{id}/default` | 会员 | 设默认地址 | — | true | — |

## 收藏 / 足迹

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/user/favorites` | 会员 | 收藏列表（实时价态） | page | list[spuId,spuName,image,price,sellable,invalid], total | 10003 |
| PUT | `/user/favorites/{spuId}` | 会员 | 收藏（已软删则复活） | — | true | — |
| DELETE | `/user/favorites/{spuId}` | 会员 | 取消收藏（软删） | — | true | — |
| GET | `/user/footprints` | 会员 | 足迹列表（按最近浏览倒序） | page | list[spuId,spuName,image,price,lastViewAt], total | 10003 |
| DELETE | `/user/footprints` | 会员 | 清空足迹 | — | true | — |

## 优惠券

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/user/coupons/available` | 会员 | 可领模板列表 | page | list[couponId,name,type,threshold,discount,validDesc,canReceive] | 10003 |
| POST | `/user/coupons/{couponId}/receive` | 会员 | 领券 | — | userCouponId | 50001 已领完/超限 |
| GET | `/user/coupons` | 会员 | 我的券 | status(1未用/2已用/3过期/4退回), page | list[userCouponId,name,threshold,discount,expireTime], total | 10003 |

## 消息中心

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/user/messages` | 会员 | 站内信列表 | isRead?, page | list[id,title,content摘要,bizType,bizNo,isRead,createdAt], total, unreadCount | 10003 |
| PUT | `/user/messages/{id}/read` | 会员 | 标记已读 | — | true | — |
| PUT | `/user/messages/read-all` | 会员 | 全部已读 | — | true | — |
| GET | `/user/notify-preferences` | 会员 | 通知偏好 | — | list[channel,enabled] | 10003 |
| PUT | `/user/notify-preferences` | 会员 | 设置偏好 | list[channel,enabled] | true | 10001 |

## 分销中心

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| POST | `/user/distribution/apply` | 会员 | 申请推广员 | — | status(待审核) | 60001 已申请/已是推广员 |
| GET | `/user/distribution/status` | 会员 | 推广员状态与等级 | — | status, level | 10003 |
| GET | `/user/distribution/relations` | 会员 | 我的邀请关系（直接上级+下级列表） | page(下级) | inviter{nickname,avatar(脱敏)}, invitees[], total | 10003 |
| GET | `/user/distribution/commission-rules` | 会员 | 可见佣金比例（按商品/分类查询） | spuId? | list[scope,level1Rate,level2Rate] | 10003 |
| GET | `/user/distribution/records` | 会员 | 佣金记录 | status?, page | list[orderNo,level,amount,status,settleTime], total | 10003 |
| GET | `/user/distribution/account` | 会员 | 佣金账户 | — | balance, frozen | 10003 |
| GET | `/user/distribution/account/logs` | 会员 | 账户流水 | bizType?, page | list[amount,balanceAfter,frozenAfter,bizType,bizNo,createdAt], total | 10003 |
| POST | `/user/distribution/withdraws` | 会员 | 提现申请 | amount | withdrawNo | 60002 余额不足, 60003 进行中存在单 |
| GET | `/user/distribution/withdraws` | 会员 | 提现列表 | status?, page | list[withdrawNo,amount,status,createdAt], total | 10003 |
| GET | `/user/distribution/invite-records` | 会员 | 邀请激励记录 | page | list[newUser(脱敏昵称),rewardType,status,createdAt], total | 10003 |
| GET | `/user/share-code` | 会员 | 我的推广码 | — | shareCode, shareLink | 10003 |
