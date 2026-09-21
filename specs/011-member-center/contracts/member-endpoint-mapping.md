# Contract: 011-member-center 端点 × 方法 × 越权（21 端点）

> 全部端点：**要求登录**（Bearer）+ **仅限本人数据**（userId 来自会话，不接受入参指定他人）。
> 无权限点（会员自助能力，非管理面）；本批端点均不在公开白名单。
> 调用形态：controller 直调 service 包级函数（user 域既有形态, research D1）。

## 资料与登录记录（3）
| 端点 | 方法+路径 | service | 越权 |
|---|---|---|---|
| ProfileDetail | GET /user/profile | `ProfileDetail(ctx, userId)` | 本人 |
| ProfileUpdate | PUT /user/profile | `ProfileUpdate(ctx, userId, in)` | 本人 |
| LoginLogList | GET /user/login-logs | `LoginLogs(ctx, userId, page)` | 本人 |

## 收货地址（5）
| 端点 | 方法+路径 | service | 越权 |
|---|---|---|---|
| AddressList | GET /user/addresses | `AddressList(ctx, userId)` | 本人 |
| AddressCreate | POST /user/addresses | `AddressCreate(ctx, userId, in)` | 本人 |
| AddressUpdate | PUT /user/addresses/{id} | `AddressUpdate(ctx, userId, id, in)` | **归属校验**（他人→10006） |
| AddressDelete | DELETE /user/addresses/{id} | `AddressDelete(ctx, userId, id)` | **归属校验** |
| AddressSetDefault | PUT /user/addresses/{id}/default | `AddressSetDefault(ctx, userId, id)` | **归属校验** + 同事务清其他 |

## 收藏与足迹（5）
| 端点 | 方法+路径 | service | 说明 |
|---|---|---|---|
| FavoriteList | GET /user/favorites | `FavoriteList(ctx, userId, page)` | 含实时价态 |
| FavoriteAdd | PUT /user/favorites/{spuId} | `FavoriteAdd(ctx, userId, spuId)` | **复活语义**（幂等） |
| FavoriteRemove | DELETE /user/favorites/{spuId} | `FavoriteRemove(ctx, userId, spuId)` | 软删（幂等） |
| FootprintList | GET /user/footprints | `FootprintList(ctx, userId, page)` | 时间倒序+浏览次数 |
| FootprintClear | DELETE /user/footprints | `FootprintClear(ctx, userId)` | 清空本人 |

## 站内信与通知偏好（5）
| 端点 | 方法+路径 | service | 说明 |
|---|---|---|---|
| MessageList | GET /user/messages | `Messages(ctx, userId, isRead, page)` | 含**未读计数**（全量, 非当前页） |
| MessageRead | PUT /user/messages/{id}/read | `MarkRead(ctx, userId, id)` | 归属校验 |
| MessageReadAll | PUT /user/messages/read-all | `MarkAllRead(ctx, userId)` | 本人全部 |
| NotifyPreferenceGet | GET /user/notify-preferences | `Preferences(ctx, userId)` | 两渠道；未设置=全开 |
| NotifyPreferenceSet | PUT /user/notify-preferences | `SetPreferences(ctx, userId, prefs)` | 幂等 UPSERT |

## 积分（2）
| 端点 | 方法+路径 | service | 说明 |
|---|---|---|---|
| PointAccount | GET /user/point-account | `PointAccount(ctx, userId)` | 余额可负（如实） |
| PointLogList | GET /user/point-logs | `PointLogs(ctx, userId, bizType, page)` | 类型筛选+分页 |

## 邀请记录（1）
| 端点 | 方法+路径 | service | 说明 |
|---|---|---|---|
| InviteRecordList | GET /user/distribution/invite-records | `InviteRecords(ctx, userId, page)` | 仅本人；被邀请人脱敏；**接口微扩** |

## 内部方法（实现, 不建端点）
`AddressGetForOrder` / `FootprintRecord` / `FootprintCleanExpired` / `NotifyEnqueue` / `NotifyDispatchTask` /
`PointEarn` / `PointConsume` / `PointRefund` / `PointExpireDormant` / `ProfileGrowthAdd` / `ProfileLevelRecalc`
——供交易链路与定时任务调用；本批不改 006 既有调用方式。
