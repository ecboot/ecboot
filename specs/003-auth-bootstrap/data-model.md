# Data Model: 认证引导纵切片（Phase 1）

**零新表**——本切片的数据形态 = 既有表的读写映射 + Redis 键模型。

## 一、Redis 键模型（infra-core 组件的"表"）

| 键 | 值 | TTL | 语义 |
|---|---|---|---|
| `captcha:img:{ticket}` | 小写答案（4 字符） | 图形有效期（默认 300s，配置 `captcha.image.ttl-seconds`） | 图形验证码凭证；校验 `GETDEL`（原子一次性） |
| `captcha:sms:{phone}` | 6 位数字码 | 短信有效期（默认 300s） | 短信验证码；校验 `GETDEL` |
| `captcha:sms:interval:{phone}` | 1 | 60s | 重发间隔锁（SET NX EX） |
| `captcha:fail:{phone}` | 计数 | 1800s | 连续失败计数（INCR；≥5 触发锁定=键存在期内拒绝） |
| `session:{token}` | userId | 7d（默认，配置 `session.ttl-days`） | 会话；校验时 EXPIRE 续期（滑动）；登出 DEL |
| `mock:sms:{phone}` | 码 | 同短信 | **仅 mock 模式**：联调取码通道 |

键前缀集中在 `SessionKeys`/常量类；全部 TTL 走 `system_config` 映射（应用层缓存 ≤30 分钟刷新）。

## 二、JPA 实体映射（service-user.domain.model → 既有表）

**User** → `user` 表（V1+V11+V19+V27+V28）：

| 字段 | 映射 | 切片行为 |
|---|---|---|
| id / nickname / avatar / gender | 直映 | 注册时昵称默认 `用户{phone后4位明文内存态}` |
| phone / phone_hash | 密文/哈希 | 写入经 `PhoneCipher`；**实体上永不出现明文字段**（命令对象持有明文，落库前转换） |
| wx_openid / wx_unionid | 直映 | 归并时写入；冲突（非空且不等）→ 业务异常 20004 |
| status | 直映 | 登录校验（2=禁用 → 20003） |
| register_channel | 直映 | 注册时按渠道写入 |
| last_login_at / last_active_at | 直映 | 登录成功同事务更新（休眠分级依据） |
| share_code | 直映 | 注册时生成（V28 规则：人人可分享），16 位码 |

**UserLoginLog** → `user_login_log`：每次登录尝试一行（user_id/login_channel/login_status/ip/user_agent）；失败也记录（user_id 可为 0=账号不存在场景，应用层约定）。

## 三、状态与规则（应用层，落 AuthService）

**注册即登录决策树**：

```text
短信码校验通过 → phone_hash 查询
 ├─ 命中 → 校验 status（禁用→20003）→ 休眠分级（last_active_at ≥ 阈值*且非本次完整核验→拒绝要求重走）
 │         → 登录成功（更新 last_*，写日志，发 token）
 └─ 未命中 → 创建（phone=密文, hash, channel, share_code）→ 首次登录即成功
```

**微信归并决策树**：

```text
WxClient 取 openid → openid 查 user
 ├─ 命中 → 直接登录（休眠分级同上）
 └─ 未命中 → 取手机号（mock 为命令字段）
      ├─ phone_hash 命中 → 目标账号 wx_openid 为空？→ 绑定+登录 ： 已绑其他 → 20004（人工渠道）
      └─ 未命中 → 新建（含 openid）→ 登录
```

**休眠分级（本切片仅第一档）**：≥ 阈值（默认 90 天，`system_config: dormant.tier1.days`——新增配置项）→ 强制完整验证码核验（本切片的短信码登录天然满足"强核验"；微信静默登录场景被拒绝并降级到手机号通道）。

**验证规则**：手机号 `^1[3-9]\d{9}$`（入口校验）；验证码 `\d{6}`；图形答案等值忽略大小写。

## 四、新增 system_config 项

| code | 默认 | 说明 |
|---|---|---|
| `captcha.image.ttl_seconds` | 300 | 图形码有效期 |
| `captcha.sms.ttl_seconds` | 300 | 短信码有效期 |
| `captcha.sms.resend_seconds` | 60 | 重发间隔 |
| `captcha.fail.max` | 5 | 连续失败上限 |
| `captcha.fail.lock_seconds` | 1800 | 锁定时长 |
| `session.ttl_days` | 7 | 会话有效期 |
| `dormant.tier1.days` | 90 | 休眠一级阈值 |

> 种子数据以 V32 迁移落库（本切片唯一 DDL：INSERT 配置行）。
