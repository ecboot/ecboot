# Quickstart: 认证引导纵切片验证指南（Go 版）

前提：`cd apps/server && docker compose up -d`（MySQL/Redis/ES）；`make migrate-up`（含 000032 配置种子）；`make run`（:8080）。端到端用例可直接对运行中服务执行（或由 gtest 端到端用例覆盖同序列）。

## 场景一：注册即登录全链路（SC-001）

```bash
BASE=http://localhost:8080
# 1. 图形验证码
curl -s $BASE/api/captcha/image | jq -r '.data.ticket,.data.imageBase64[:40]'
#    → 记下 ticket；图片肉眼识别答案（自动化测试中经组件直取答案）
# 2. 校验图形码（一次性：重复用同 ticket 必须失败）
curl -s -X POST $BASE/api/captcha/verify -H 'Content-Type: application/json' \
     -d '{"ticket":"<ticket>","answer":"<答案>"}'
# 3. 发短信（mock 模式）
curl -s -X POST $BASE/api/captcha/sms -H 'Content-Type: application/json' \
     -d '{"phone":"13800001111","ticket":"<ticket2>","answer":"<答案2>"}'
# 4. 取 mock 码（仅 mock 模式存在）
curl -s "$BASE/api/captcha/sms/mock-latest?phone=13800001111"
# 5. 登录（新号自动注册）
curl -s -X POST $BASE/api/auth/sms-login -H 'Content-Type: application/json' \
     -d '{"phone":"13800001111","smsCode":"<码>"}'   # → data.token, isNew=true
# 6. 受保护端点 + 登出后失效（SC: 会话生命周期）
curl -s $BASE/api/auth/me -H "Authorization: Bearer <token>"       # → 用户信息
curl -s -X POST $BASE/api/auth/logout -H "Authorization: Bearer <token>"
curl -s $BASE/api/auth/me -H "Authorization: Bearer <token>"       # → code=10003
```

**预期**：5 步拿到 token；登出后 me 返回 10003；全部响应为统一三段式结构（SC-002）。

## 场景二：防刷规则正反用例（SC-003）

| 用例 | 操作 | 预期 |
|---|---|---|
| 一次性 | 同一 ticket 校验两次 | 第二次 code=20001 |
| 有效期 | 等待 TTL 后校验 | code=20001 |
| 重发频控 | 60 秒内两连发 | 第二次 code=10004 |
| 失败锁定 | 连续 5 次错码登录 | 第 5 次后 code=20002 |

## 场景三：微信归并（SC-004）

```bash
# 先用手机号 13800002222 注册（同场景一）→ 再以微信身份携同号登录
curl -s -X POST $BASE/api/auth/wx-login -H 'Content-Type: application/json' \
     -d '{"wxCode":"dev001","phone":"13800002222"}'
# 预期：登录成功且 userId 与手机号账号一致（未新建账号）
# 再以 wxCode=dev001 直接登录 → 成功（openid 已绑定）
# wxCode=dev002 + 同手机号 → code=20004（绑定冲突）
```

DB 侧断言：`SELECT COUNT(*) FROM user WHERE phone_hash IS NOT NULL` 与账号数一致；**`user` 表 phone 列无明文**（SC-007，`LENGTH(phone)<>64 OR phone REGEXP '^1[3-9]'` 计数为 0）；`user_login_log` 每次尝试一行（SC-005）。

## 场景四：休眠分级（SC-006）

```sql
-- 手工将账号活跃时间拨回 91 天前
UPDATE user SET last_active_at = DATE_SUB(NOW(), INTERVAL 91 DAY) WHERE id = <uid>;
```

再次微信静默登录该账号 → 被拒并引导手机号验证码通道（短信码登录可用）。

## 宪法 IV 验证命令

```bash
cd apps/server && make test    # go test ./... 全绿
```
