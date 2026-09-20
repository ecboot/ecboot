# common 渠道契约（公共）

前缀 `/common`。除标注外全部公开。

## 验证码（按工程实现回记, 2026-09-20——消除与既有代码的双契约）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/common/captcha` | 公开 | 获取图形验证码 | — | captchaKey, captchaImg | — |
| POST | `/common/sms/code` | 公开 | 发送短信验证码（先过图形码） | phoneNumber, template, captchaCode, captchaKey | smsCodeKey | 10004, 20001, 20002 |
| POST | `/common/captcha/verify` | 公开 | 图形码独立校验（预留, 随 003 实现补） | ticket/Key, answer | true | 20001 |
| GET | `/common/captcha/sms/mock-latest` | 公开（仅 mock 模式注册） | 联调取码 | phone | smsCode | — |

> 实现语义以既有代码（api/common/v1/captcha.go、sms.go）为准；verify 与 mock-latest 端点随 003 实现补齐。

## 门店（游客浏览，store 域）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/common/stores` | 公开 | 门店列表（区县筛选或附近检索） | districtCode?; longitude?,latitude?,radiusKm?, page | list[storeNo,name,address,distanceM?,businessHours,pickupEnabled,status], total | — |
| GET | `/common/stores/{id}` | 公开 | 门店详情 | — | 全字段（含经纬度、电话） | 70001 门店不存在/已歇业提示语义 |

## 公共

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/common/ping` | 公开 | 健康探针 | — | pong | — |
| POST | `/common/shares` | 公开（可选会员身份） | 分享行为上报（归因来源） | spuId?, shareChannel, sceneValue? | true | — |
