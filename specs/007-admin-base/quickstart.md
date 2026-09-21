# Quickstart: 007-admin-base 验证指南

> 目标：不读实现代码即可验证批次 DoD。实现细节见 [plan.md](./plan.md) 与 tasks.md（后续生成）。

## 前置

```bash
cd apps/server
docker compose up -d          # MySQL 13306 / Redis 6379
make migrate-up               # 应用迁移至 000035（含种子超管）
```

## 自动化验证（DoD 主通道）

```bash
make test                     # 全绿（含既有认证/交易回归 + 本批 service/system 新测试）
make lint                     # 通过
make check-stub               # admin 桩 129→109、common 桩 5→3（本批 22 清零）
```

测试基座遵循 005/006 模式：`internal/service/system/*_test.go` 确定性配置注入 + 数据自建清理。

## 手工验证序列（可选，curl）

种子超管账号见 `migrations/000035_admin_seed.up.sql`（首次登录后应改密）。

1. **登录**：`POST /admin/login`（username/password）→ 返回 token/refreshToken/realName/isSuper=true；
   错误密码 → 统一错误码；`admin_login_log` 新增对应成功/失败行。
2. **个人信息**：`GET /admin/profile`（Bearer token）→ username/realName/roles。
3. **权限拦截**：创建非超管账号 A（`POST /admin/admin-users`）→ A 登录 → A 调 `POST /admin/roles` →
   权限不足；超管执行同一请求 → 成功。
4. **角色分配**：建角色 → `PUT /admin/roles/{id}/permissions` 全量替换（重复提交幂等）→
   角色详情 permissionIds 与提交一致。
5. **账号治理**：禁用 A → A 旧 token 调 `GET /admin/profile` 被拒；软删 A → 列表不可见。
6. **配置覆盖层**：`GET /admin/configs` 见种子配置 → `PUT /admin/configs/{code}` 改值 → 列表反映新值；
   置 status=0 → `session.ttl_days` 类消费方回退默认（既有 FromConfig 行为）。
7. **会话生命周期**：登出（双凭证同毁）→ 原 token 被拒；refreshToken 重放被拒。
8. **基础端点**：`GET /common/ping` 无凭证 200；`GET /common/mock-latest-sms` 返回最近模拟验证码
   （非生产环境）。

## 验收红线（宪法 IV）

任何"完成"声明前，`make test` 与 `make lint` 的实际输出必须为通过；桩数以 `make check-stub` 输出为准。
