# ecboot-api-user（会员中心渠道）

小程序前台·**会员中心** REST API——围绕"我是谁、我的资产、我的关系"的会员侧端点。

## 功能内容（REST 端点规划）

| 功能组 | 端点（规划） | 后端域 |
|---|---|---|
| 账号 | 注册（手机号+短信验证码）、登录（手机号/微信双通道）、登出、**注销**（二次确认+阻断校验+匿名化） | service-user · 账号 |
| 个人资料 | 昵称/头像查询与修改、等级与成长值展示 | service-user |
| 收货地址 | 地址 CRUD、设默认 | service-user |
| 收藏与足迹 | 收藏/取消（复活语义）、收藏列表（实时价态）、足迹查看/清空 | service-user |
| 我的资产 | 积分账户与流水、优惠券（领券/可用/已用/已过期） | service-user |
| 消息中心 | 站内信列表/已读、订阅偏好 | service-user |
| 分销中心 | 推广员申请、我的关系链、佣金记录与账户、提现申请、邀请激励记录 | service-user · 分销 |
| 登录记录 | 近 30 天登录日志（安全中心） | service-user |
| 验证码前置 | 发码/图形码经 `api-common` 渠道端点获取（**运行时 HTTP 并存，非编译依赖**） | api-common 渠道 |

## 职责边界

- **做**：会员侧 REST 暴露（薄控制器 + DTO 组装 + 参数校验），编排 service-user
- **不做**：与 api-common/shop/admin **互不编译依赖**（enforcer 四渠道互禁）；不含商城交易端点（在 api-shop）；不含后台管理端点（在 api-admin）

## 依赖关系

- 白名单：`ecboot-api-webmvc`、`ecboot-service-user`、`ecboot-service-shop`、`ecboot-common`、`ecboot-infra-core`
- 被依赖：仅 `ecboot-start`

相关文档：`../README.md`、`ecboot-service-user/README.md`
