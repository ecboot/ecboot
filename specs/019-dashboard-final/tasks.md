---
description: "任务清单：看板与收口终验（批次 13）"
---

# Tasks: 看板与收口终验（019-dashboard-final）

**Feature**: `specs/019-dashboard-final` | **Spec**: [spec.md](./spec.md)

**Input**: 3 admin 端点；实现 `IDashboardLogic`（接口/DTO 已有）；**全量收口终验**。零迁移零权限种子。

- [ ] T001 测试（红）: 三看板已知数据口径（FR-1~3）+ 空窗口全零 + 不存在数据不报错
- [ ] T002 实现 `dashboard_impl.go`（Trade/Member/Product, 口径复用各域权威查询）; 连线 3 桩
- [ ] T003 端点级测试: dashboard:read 挂载（非超管 10005 对照——批次 10 I4 内化）
- [ ] T004 `go test ./...` 两连跑全绿; `golangci-lint run` 0 issues
- [ ] T005 `check-stub` 对账: **四渠道 3→0, 210 桩全部清零（SC-1 终验）**
- [ ] T006 收口终验清单: schema_migrations=43 与迁移文件一致/全批台账 ✅/批次 12 修复轮确认
- [ ] T007 更新 `specs/PROGRESS.md` 批次 13 状态 ✅ 与完成 commit（同 commit）并提交

## Notes
- 只读统计零资金风险; 口径正确性=全部价值（复用各域既有查询）
- 防线铁律内化（*int 不适用——无写方法; RequirePerm 首行+非超管对照）
