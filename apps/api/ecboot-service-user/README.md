# ecboot-service-user（用户领域服务）

**会员域**的领域实现——账号、资产与会员关系的业务中枢。下一步按领域拆分、内部 **DDD 分层**实现（模块边界与对外契约不变）。

## 包结构约定（DDD 四层）

```
org.juling.ecboot.user
├── interfaces/        # 对外暴露：REST DTO 组装接口 + 领域事件订阅入口（供渠道层调用）
├── application/       # 用例层：应用服务（事务边界、编排、幂等）
├── domain/            # 领域层：聚合根/实体/值对象/领域服务/仓储接口 + 领域事件定义
└── infrastructure/    # 基础设施：JPA 仓储实现、缓存/第三方适配（调 infra-core 组件）
```

## 功能内容（对齐 Schema 表）

| 功能域 | 内容 | 对应表 |
|---|---|---|
| 账号与认证 | 手机号+微信双通道注册登录、密码哈希、账号禁用、**注销匿名化**（墓碑哈希，合规红线）；**登录时手机号优先归并**——微信登录经 getPhoneNumber 取号，`phone_hash` 命中旧账号则绑定 openid 直接登录（双账号冲突消解在入口，不做迁移式合并，历史双号走后台人工改绑+审计） | `user`（V1/V11：phone 密文+phone_hash） |
| 登录审计 | 每次登录结果留痕、用户侧近 30 天查询 | `user_login_log` |
| 收货地址 | 地址 CRUD、默认地址切换、快照供下单 | `user_address` |
| 收藏与足迹 | 收藏（唯一/复活语义）、足迹（去重+90 天清理） | `user_favorite`、`user_footprint` |
| 积分账本 | 签到/消费/分享获取、下单消耗、退款回退（可负余额）、全部流水 | `point_account`、`point_log` |
| 成长值与等级 | 成长值只增不减、等级门槛判定与权益（双账本分离） | `user.growth_value/level`、`user_level_rule` |
| 优惠券持有 | 领券（防超发）、核销/退回、可用性查询（模板在 shop 域） | `user_coupon` |
| 消息中心 | 站内信必达、通知任务多渠道投递（重试上限）、订阅偏好 | `user_message`、`notify_task` |
| 分销（合规核心） | **两级封顶关系链**（结构强制，ADR-0003）、**分享归因**（分享>关系链>自然流量三级佣金判定，窗口默认 7 天）、推广员资质、佣金计提与冲销、账户（可负）、提现（渠道单号幂等）、邀请激励（注册/首单双时机） | `user_relation`、`share_record`、`distribution_user`、`commission_rule/record`、`user_account`、`account_log`、`withdraw_order`、`invite_record` |

## 职责边界

- **做**：会员域全部业务规则；领域事件发布（如"佣金结算完成"供 shop 域订阅）
- **不做**：不暴露 REST（渠道层职责）；不依赖 api 渠道/兄弟服务/start（enforcer）；商品与交易在 shop 域（券模板/订单由 shop 拥有，本域只持有券实例）

## 依赖关系

- 白名单：`ecboot-common`、`ecboot-infra-core`（兄弟 service-shop 互禁，经领域事件协作）
- 被依赖：api-user / api-admin 渠道

相关文档：ADR-0003（两级封顶）、`docs/schema-design.md` §注销匿名化
