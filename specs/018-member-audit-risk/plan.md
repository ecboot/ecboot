# Implementation Plan: 会员管理与审计风控（批次12）

**Branch**: `018-member-audit-risk` | **Spec**: [spec.md](./spec.md) | **Created**: 2026-09-22

## Technical Context

- 13 admin 端点；`IAuditLogic`（system 域）已有待实现；**新建** `IMemberAdminLogic`（user 域 4 方法）与 `IRiskAdminLogic`（system 域 7 方法, 风控属平台治理域）+ 风控评估器（`shop.IRiskHit` 端口实现, bootstrap 装配）。
- 表全既有：user / risk_rule / risk_record（000021）/ admin_login_log / admin_operation_log。**零迁移**；权限 7 点全就位（000032）。

## Constitution Check

| 门 | 判定 |
|---|---|
| 模块化单体/分层 | member 管理落 user 域（数据归属）; 风控管理/审计落 system 域（平台治理）; 端口实现落 system, bootstrap 装配 |
| 红线 | 无涉（无分销/资金面） |
| 可验证交付 | SC-1 桩数 16→3; SC-3 风控全链一条龙测试 |
| 简单优先 | 评估器 V1 黑名单全量+记录面完整; 复杂实时统计判定记账 |

## Design Decisions

- **D1 防线铁律内化**: 禁用/申诉/规则状态变更=条件 UPDATE 判行数; 改绑手机号 uk_phone_hash 1062→业务码; 未传语义 *int 三态（批次 10 I6 教训）。
- **D2 评估器**: `RiskHit.Hit(ctx, userId, ruleType, payload)` → 查启用规则（type+condition_expr 含用户标识）→ 命中: 落 risk_record（按规则 action 1拦截/2标记）+ 返回拦截; 失败/无规则 → 放行（批次 09 降级语义延续, 装配后为真实判定）。
- **D3 接口新建（记账）**: `IMemberAdminLogic`（user 域: List/Detail/Disable/RebindPhone）; `IRiskAdminLogic`（system 域: RuleList/Create/Update/Delete/Detail、RecordList、Appeal）。
- **D4 脱敏**: 列表手机号沿用既有掩码口径（批次 05 I1: 脱敏仅会员边界——后台列表仍脱敏展示）。
- **D5 申诉语义**: appeal_status（表注释口径）条件更新, 已有结论拒绝; 结论必填。

## Risks

- 风控评估器是**新写入路径**（risk_record 落库在业务主流程内）——评估失败不阻断主流程 + 落库失败仅告警。
