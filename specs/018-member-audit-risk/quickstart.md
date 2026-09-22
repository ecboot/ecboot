# Quickstart: 会员管理与审计风控（018）冒烟

**Prerequisites**: 应用库 `ecboot` 迁移版本 43（本批零迁移）。

1. **会员治理**: 列表可见（手机号脱敏）→ 禁用 → 会员登录被拒 → 启用恢复 → 改绑新手机号 → 详情一致 → 新号被占 → 业务码拒绝。
2. **风控全链**: 建黑名单规则（指向用户 U）→ U 触发玩法入口 → 被拦截且 risk_record 落行(action=1) → 申诉通过 → 记录更新; 停用规则 → 放行。
3. **审计**: 批次 01 登录产出的 admin_login_log 行 → 列表分页可见; 操作日志同。

```bash
cd apps/server && go test ./... -count=1 && golangci-lint run
for ch in admin common shop user; do echo "$ch stub=$(grep -l CodeNotImplemented internal/controller/$ch/*.go | wc -l)"; done  # admin 3 / 0 / 0 / 0
```
