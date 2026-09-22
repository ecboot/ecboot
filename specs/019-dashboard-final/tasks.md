---
description: "任务清单：看板与收口终验（批次 13）"
---

# Tasks: 看板与收口终验（019-dashboard-final）

**Feature**: `specs/019-dashboard-final` | **Spec**: [spec.md](./spec.md)

**Input**: 3 admin 端点；实现 `IDashboardLogic`（接口/DTO 已有）；**全量收口终验**。零迁移零权限种子。

- [x] T001 测试（红）: 三看板已知数据口径（FR-1~3）+ 空窗口全零 + 不存在数据不报错
- [x] T002 实现 `dashboard_impl.go`（Trade/Member/Product, 口径复用各域权威查询）; 连线 3 桩
- [x] T003 端点级测试: dashboard:read 挂载（非超管 10005 对照——批次 10 I4 内化）
- [x] T004 `go test ./...` 两连跑全绿; `golangci-lint run` 0 issues
- [x] T005 `check-stub` 对账: **四渠道 3→0, 210 桩全部清零（SC-1 终验）**
- [x] T006 收口终验清单: schema_migrations=43 与迁移文件一致/全批台账 ✅/批次 12 修复轮确认
- [x] T007 更新 `specs/PROGRESS.md` 批次 13 状态 ✅ 与完成 commit（同 commit）并提交

## Notes
- 只读统计零资金风险; 口径正确性=全部价值（复用各域既有查询）
- 防线铁律内化（*int 不适用——无写方法; RequirePerm 首行+非超管对照）

## 完成记录（2026-09-22）

- **3 端点全清 → 210 桩全部清零**: `check-stub` 四渠道 **0/0/0/0**; controller 文件 218 = api 端点定义 218（admin 129/common 8/shop 42/user 39）与 §二 基线快照一致
- **实现**: `dashboard_impl.go`（Trade/Member/Product 三方法, 口径全部锚定既有域查询: 已支付 status>=20 销售额 pay_amount 合计/退款额售后完成 status=50 实退/休眠 last_active_at≥**配置阈值**（读 system_config `dormant.tier1.days`, 与 user/wx.go 同源; 休眠定义在 **000027**）/低库存 available<=warn_count 复用 inventory 预警口径/待审核评价 audit_status=0）; 连线 3 桩
- **实现期新发现**: `Sum()` 返回 float64 违反"禁 float 存算金额"铁律 → 改 `Value("COALESCE(SUM(...),0)")` + money 规范化到分
- **收口终验**: 迁移账本一致（schema_migrations=43 = migrations/**.up.sql 43 个, dirty=0）; 全量两连跑全绿; lint 0; 13 批全部 ✅
- **端点级测试**: 3 看板超管可达 + 非超管 10005 对照（dashboard:read 挂载守卫）
- **验证**: `go test ./...` 两连跑全绿; `golangci-lint` **0 issues**; 零迁移零权限种子
