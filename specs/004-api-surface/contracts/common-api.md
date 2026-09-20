# common 渠道契约（公共）

前缀 `/common`。除标注外全部公开。

## 验证码（003 已设计，路径迁移版）

| 方法 | 路径 | 鉴权 | 说明 | 关键入参 | 关键出参 | 错误码 |
|---|---|---|---|---|---|---|
| GET | `/common/captcha` | 公开 | 获取图形验证码 | — | ticket, imageBase64, expiresIn | — |
| POST | `/common/captcha/verify` | 公开 | 图形码一次性校验 | ticket, answer | true | 20001 |
| POST | `/common/captcha/sms` | 公开 | 发送短信验证码（先图形码） | phone, ticket, answer | expiresIn | 10004, 20001, 20002 |
| GET | `/common/captcha/sms/mock-latest` | 公开（仅 mock 模式注册） | 联调取码 | phone | smsCode | — |

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
