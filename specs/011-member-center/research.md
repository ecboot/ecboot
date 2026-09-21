# Research: 011-member-center

> Phase 0 产出。基于代码勘察（commit 4ba9422 时点）。

## D1 实现形态 = 包级函数（跟随 user 域既有）

**结论**：新增实现一律用包级函数（`func ProfileDetail(ctx, userId) (...)`），与 `SmsLogin`/`WxLogin` 先例一致；
controller 直接调用。

**理由**：user 域既有实现（auth/wx）为包级函数；同域单一形态。批次 04 的 shop 域连线用 struct 是因既有
`ProductLogicImpl` 如此——形态跟随**域内既有**，不跨域统一。

## D2 通知偏好表设计（新增迁移 000036）

**结论**：
```sql
CREATE TABLE `user_notify_preference` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `channel` TINYINT NOT NULL COMMENT '1小程序订阅 2短信',
  `enabled` TINYINT NOT NULL DEFAULT 1,
  `created_at`/`updated_at`,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_channel` (`user_id`,`channel`)
)
```

**理由**：勘察确认偏好**无任何既有存储**（无表无列）；接口语义（Enqueue 按偏好拆渠道）要求持久化。
唯一键 (user_id, channel) 支撑幂等 UPSERT；**未设置 = 不建行 = 默认全开**（避免预置行，YAGNI）。
站内信不在偏好内（api 的 NotifyPreference 仅两渠道, 与接口注释"站内信除外"一致）。

## D3 手机号脱敏（服务层）

**结论**：复用既有 `PhoneCipherFor(ctx)` 解密 + 掩码（保留前 3 后 4, 中间 `****`）；
`ProfileDetail.MaskedPhone` 为 DTO 字段（明文不落 DTO）。

**理由**：延续认证纵切片的加密存储与脱敏约束（个保法）。

## D4 越权防护收口

**结论**：所有会员自助方法与端点均以 `userId`（来自 `middleware.CtxUserIdFrom`）为第一约束——
查询带 `Where(user_id, userId)`，改删带 `Where(id, x).Where(user_id, userId)`（或先校验归属再操作），
他人资源一律"不可见/不存在"语义。

## D5 invite_record：接口微扩 + DTO 新增

**结论**：`IDistributionLogic` 增 `InviteRecords(ctx, userId, page)`；`model` 新增 `InviteRecordItem`
（对齐 api：NewUser 脱敏昵称/RewardDesc/Status/CreatedAt）。

**理由**：勘察确认该接口无此方法、model 无此 DTO；api 契约已定义端点（011 清单归属本批）。

## D6 内部方法：实现但不暴露端点

**结论**：`Footprint.Record`/`Address.GetForOrder`/`Notify.Enqueue·DispatchTask`/`Point.Earn·Consume·Refund·ExpireDormant`/
`Profile.GrowthAdd·LevelRecalc` 一并实现（接口已定义），本批不建端点。

**理由**：它们是交易链路与定时任务的调用面；006 交易链路当前以自实现方式处理（如库存直操作），
本批**不改其调用**（避免回归），仅把能力补齐供后续批次切换。

## 勘察结论（非决策）

- **DTO 齐备**：`dto_user.go` 23 类型含本批绝大多数出入参（Address/Favorite/Footprint/Message/NotifyPreference/
  PointAccountView/PointLogItem/ProfileDetail/ProfileUpdateInput/UserLoginLogItem）。
- **表齐备**（除偏好）：user_address/user_favorite/user_footprint/user_message/notify_task/notify_template/
  point_account/point_log/user_level_rule/user_login_log/invite_record。
- **user 扩展列**：昵称/头像/性别（000001）+ 成长值等（V11/V19/V27/V28 ALTER）。
- **测试基座**：`internal/service/user/auth_flow_test.go` 的 init（确定性 gdb/gredis + UTC 时区）+ `issueCode`/
  `cleanState` 惯例可复用。
