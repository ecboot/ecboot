# 端点到实现映射：会员管理与审计风控（018）

**13 端点**（验收：admin 16→3, 四渠道 16→3）。

## 接口
- 新建 `IMemberAdminLogic`（user 域）: AdminList/AdminDetail/Disable/RebindPhone（D3）
- 新建 `IRiskAdminLogic`（system 域）: RuleList/Create/Update/Delete/Detail/RecordList/Appeal（D3）
- 实现 `IAuditLogic`（system 域, 已有）: OperationLogs/LoginLogs
- 风控评估器: 实现 `shop.IRiskHit`（端口, bootstrap 装配——批次 09/11 降级点闭合, FR-6）

## 端点映射

| # | 端点 | controller 桩 | service | 权限 |
|---|---|---|---|---|
| 1 | `GET /admin/members` | admin_v1_admin_member_list.go | MemberAdmin.AdminList | member:read |
| 2 | `GET /admin/members/{userId}` | admin_v1_admin_member_detail.go | .AdminDetail | member:read |
| 3 | `POST /admin/members/{userId}/disable` | admin_v1_admin_member_disable.go | .Disable | member:update |
| 4 | `POST /admin/members/{userId}/rebind-phone` | admin_v1_admin_member_rebind_phone.go | .RebindPhone | member:update |
| 5-8 | `GET/POST/PUT/DELETE /admin/risk-rules`(+`/{id}`) | admin_v1_admin_risk_rule_*.go | RiskAdmin.Rule{List,Create,Update,Delete} | risk:rule:read/manage |
| 9 | `GET /admin/risk-rules/{id}`(详情) | admin_v1_admin_risk_rule_detail.go | .RuleDetail | risk:rule:read |
| 10 | `GET /admin/risk-records` | admin_v1_admin_risk_record_list.go | .RecordList | risk:record:read |
| 11 | `POST /admin/risk-records/{id}/appeal` | admin_v1_admin_risk_appeal.go | .Appeal | risk:record:appeal |
| 12 | `GET /admin/login-logs` | admin_v1_admin_login_log_list.go | Audit.LoginLogs | system:audit:read |
| 13 | `GET /admin/operation-logs` | admin_v1_admin_operation_log_list.go | Audit.OperationLogs | system:audit:read |

## 校验矩阵
- 改绑手机号: 新号格式+`uk_phone_hash` 1062→业务码; 条件更新判行数
- 禁用: 条件更新（status 期望态判行数）; 同值幂等
- 规则: type∈1..5 / action∈1,2 / condition_expr 非空; 软删
- 申诉: appeal_status 条件更新（已结论拒绝）; 结论必填
