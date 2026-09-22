---
description: "任务清单：会员管理与审计风控（批次 12）"
---

# Tasks: 会员管理与审计风控（018-member-audit-risk）

**Feature**: `specs/018-member-audit-risk` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/member-risk-endpoint-mapping.md](./contracts/member-risk-endpoint-mapping.md)

**Input**: 13 admin 端点；新建 IMemberAdminLogic（user 域 4 方法）+ IRiskAdminLogic（system 域 7 方法）+ 实现 IAuditLogic；风控评估器（IRiskHit 装配）。零迁移零权限种子。

## Phase 1: Foundational

- [ ] T001 接口新建（D3 记账）: `IMemberAdminLogic`（user 域）/`IRiskAdminLogic`（system 域）; DTO 勘察补缺
- [ ] T002 [P] 测试基座: 会员/风控规则/记录 fixture 与自清

## Phase 2: US1 会员管理（P1）🎯 MVP

- [ ] T003 [P] [US1] 测试（红）: 列表筛选分页+手机号脱敏/详情/**禁用条件更新+既有登录口径拒绝**/启用/**改绑唯一键 1062→业务码**/不存在 404 类
- [ ] T004 [US1] 实现 `member_admin_impl.go`（user 域 4 方法）; 连线 4 桩

## Phase 3: US2 风控规则/事件/申诉 + 评估器（P1）

- [ ] T005 [P] [US2] 测试（红）: 规则 CRUD+详情/校验矩阵/软删/**评估器: 启用黑名单命中拦截+落记录、停用不生效、无规则放行、评估失败不阻断**/申诉条件更新（已结论拒绝）
- [ ] T006 [US2] 实现 `risk_admin_impl.go`（system 域 7 方法）+ 风控评估器（实现 shop.IRiskHit）; 连线 7 桩; bootstrap 装配评估器
- [ ] T007 [US2] 风控全链一条龙测试（SC-3: 规则→命中拦截→落记录→申诉→放行）

## Phase 4: US3 审计日志（P2）

- [ ] T008 [P] [US3] 测试（红）: 登录日志/操作日志分页与筛选（复用批次 01 产出数据）
- [ ] T009 [US3] 实现 `IAuditLogic.OperationLogs/LoginLogs`; 连线 2 桩

## Phase 5: Polish & 收尾

- [ ] T010 端点级测试: **非超管 10005 对照**（全权限点）+ 超管可达 + 装配断言（评估器已注入 RiskHit）
- [ ] T011 `go test ./...` 两连跑全绿; `golangci-lint run` 0 issues
- [ ] T012 `check-stub` 对账: admin 16→3（四渠道 16→3）; quickstart 冒烟
- [ ] T013 更新 `specs/PROGRESS.md` 批次 12 状态 ✅ 与完成 commit（同 commit）并提交

## Notes（本批硬约束）

- 防线铁律内化（条件更新判行数/1062 转业务码/*int 三态/同值幂等）
- 评估器失败**不阻断主流程**（降级放行+告警, 批次 09 语义延续但装配后为真实判定）
- 风控落库在业务主流程内: 落库失败仅告警, 不影响业务结果
- 禁止手改生成物; 实时统计判定（高频/套利）记账增强
