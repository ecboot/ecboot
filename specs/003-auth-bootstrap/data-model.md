# Data Model: 认证引导纵切片（Phase 1，Go 版）

**表结构零变更**——数据形态 = 既有表的读写（经 gf gen dao）+ Redis 键模型 + 1 个配置种子迁移。

## 一、Redis 键模型（gredis，与 Java 版设计一致）

| 键 | 值 | TTL | 语义 |
|---|---|---|---|
| `captcha:img:{ticket}` | 小写答案 | 图形有效期（默认 300s） | 图形码凭证；校验 `GETDEL` 原子一次性 |
| `captcha:sms:{phone}` | 6 位码 | 短信有效期（默认 300s） | 短信码；校验 `GETDEL` |
| `captcha:sms:interval:{phone}` | 1 | 60s | 重发间隔锁（SET NX EX） |
| `captcha:fail:{phone}` | 计数 | 1800s | 连续失败计数（INCR；≥上限锁定） |
| `session:{token}` | userId | 7d | 会话；校验 EXPIRE 续期（滑动）；登出 DEL |
| `mock:sms:{phone}` | 码 | 同短信 | **仅 mock 模式**联调取码 |

键常量集中在 `internal/infra/security/keys.go`；TTL 读 `system_config`（应用内缓存 ≤30 分钟刷新）。

## 二、数据访问（gf gen dao + 领域封装）

- `make gen` 生成 `internal/dao`（user/user_login_log 等表操作）与 `internal/model`（entity/do 结构体）——生成物勿手改（宪法 IV）。
- `internal/service/user/repo.go` 在 dao 之上做领域封装：`FindByPhoneHash` / `FindByOpenid` / `CreateLogin`（密文+哈希+share_code 生成）/`TouchLoginTimes`（last_login_at/last_active_at）/`AppendLoginLog`。
- **明文手机号只存在于请求结构体与内存**：进入 repo 前必经 `PhoneCipher`（密文+哈希），entity 字段即库列（无明文字段）。

## 三、状态与规则（service/user/auth.go 决策树，语义同 Java 版）

**注册即登录**：短信码 GETDEL 校验通过 → `phone_hash` 查询 → 命中：status 校验（禁用→20003）+ 休眠分级（`last_active_at` ≥ 阈值 → 强制完整核验，本切片短信码登录天然满足；微信静默路径被拒引导手机号通道）→ 更新时间戳+写日志+发 token；未命中：创建（密文/哈希/渠道/share_code）→ 即登录。

**微信归并**：openid 查询 → 命中直接登录（休眠分级同上）；未命中 → 取手机号 → `phone_hash` 命中：目标账号 openid 为空则绑定登录，非空且不等 → 20004（人工渠道）；未命中 → 新建（含 openid）→ 登录。

**校验规则**：手机号 `^1[3-9]\d{9}$`（入口拦截，非法不进发码）；验证码 `\d{6}`；图形答案忽略大小写等值。

**并发注册**：`phone_hash` 唯一索引兜底——冲突方捕获重复键错误后转登录路径。

## 四、迁移 000032_auth_config.up.sql（本切片唯一 DDL）

`system_config` 种子 7 行：

| code | 默认 | 说明 |
|---|---|---|
| `captcha.image.ttl_seconds` | 300 | 图形码有效期 |
| `captcha.sms.ttl_seconds` | 300 | 短信码有效期 |
| `captcha.sms.resend_seconds` | 60 | 重发间隔 |
| `captcha.fail.max` | 5 | 连续失败上限 |
| `captcha.fail.lock_seconds` | 1800 | 锁定时长 |
| `session.ttl_days` | 7 | 会话有效期 |
| `dormant.tier1.days` | 90 | 休眠一级阈值 |
