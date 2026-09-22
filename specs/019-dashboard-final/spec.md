# Feature Specification: 看板与收口终验（批次13）

**Feature Branch**: `019-dashboard-final` | **Created**: 2026-09-22 | **Status**: Draft

**Input**: 批次 13（**最后一批**）：**3 个 admin 看板端点**（全为桩）——交易看板（订单数/销售额/退款额/待发货单数）、会员看板（新增/活跃/休眠（阈值可配））、商品看板（在售/低库存预警/待审核评价）；`IDashboardLogic`（Trade/Member/Product）接口与 DTO 已有（与 api 逐字段同形）；权限 `dashboard:read` 已就位。**本批同时承担全量收口终验**：210 桩清零对账、12 批全量回归、批次 12 修复轮复审确认。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 三看板（Priority: P1）🎯 MVP

运营者打开后台首页三个看板：交易看板（时间窗内订单数/销售额/退款额、当前待发货单数）、会员看板（时间窗新增/活跃/休眠（阈值可配））、商品看板（在售商品数/低库存预警数/待审核评价数）。全部为只读统计，口径复用各域既有权威查询（不自造 SQL 口径）。

**Why this priority**: 唯一剩余端点组；全部只读无资金/状态变更，风险最低——但**口径正确性**是看板的全部价值。

**Independent Test**: 造已知数据（订单 2 笔其中 1 笔已支付、会员 1 新增、在售商品 1、待审评价 1）→ 三看板数字与手工口径一致。

**Acceptance Scenarios**:

1. **Given** 时间窗内有 2 笔已支付订单（行实付合计 130.00）与 1 笔待发货，**When** 查交易看板，**Then** orderCount/salesAmount/pendingDeliver 与数据一致。
2. **Given** 时间窗内新增会员 1、90 天内活跃 1、最后活跃 ≥90 天会员 1，**When** 查会员看板，**Then** 三数一致。
3. **Given** 在售 SPU 1、库存预警 SKU 1、待审核评价 1，**When** 查商品看板，**Then** 三数一致。
4. **Given** 无数据窗口，**When** 查询，**Then** 全零而非报错。
5. **Given** 非超管（无 dashboard:read），**Then** 10005。

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-1**: 交易看板——窗口内订单数（**已支付枚举 20/30/40, 显式排除 90 已取消**——收口终验 C1 修正: cancelBy 只改状态不清零 pay_amount, `>=20` 会把取消单计成销售额）、销售额（pay_amount 合计）、退款额（售后完成 status=50 实退合计）、当前待发货单数（status=20）。
- **FR-2**: 会员看板——窗口新增（created_at）、活跃（last_active_at 窗口内; **无窗口默认近 30 天**——I4 统一口径）、休眠（阈值读 `dormant.tier1.days`（缺省 90）, 与 user/wx.go 分级实现同源; 休眠定义在 **000027**）。
- **FR-3**: 商品看板——在售 SPU 数、低库存预警数（inventory available≤warn_count 既有口径）、待审核评价数（audit_status=0）。
- **FR-4**: 权限 `dashboard:read` 挂接三端点（RequirePerm 首行）；已就位零种子。

## Success Criteria *(mandatory)*

- **SC-1**: `check-stub` 对账 **admin 3→0**——**四渠道合计 3→0，13 批 210 桩全部清零**。
- **SC-2**: 全量 `go test ./...` 两连跑全绿；`golangci-lint run` **0 issues**。
- **SC-3**: 三看板口径经已知数据测试钉住；权限挂载经非超管对照测试。
- **SC-4**: 收口终验清单完成（桩数/迁移账本/schema_migrations 与迁移文件一致/PROGRESS §三 13 批全 ✅）。

## Assumptions

- 活跃口径=last_active_at 在窗口内, 无窗口取近 30 天（I4 收口终验统一三处表述）。
- 退款额取售后单实退合计（refund_amount，status=50 已完成）；口径与售后域 status 枚举对齐。
- 时间窗为空时: 交易/会员的窗口类指标取**近 30 天**（活跃）或全量（新增/订单——按各口径注释固定）; req 字段可选。
