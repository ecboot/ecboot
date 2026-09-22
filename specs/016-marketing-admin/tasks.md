---
description: "任务清单：营销后台（批次 10）"
---

# Tasks: 营销后台（016-marketing-admin）

**Feature**: `specs/016-marketing-admin` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/marketing-admin-endpoint-mapping.md](./contracts/marketing-admin-endpoint-mapping.md)

**Input**: 30 端点（admin，全为桩）、`ICouponLogic`+`IActivityLogic` 补实现（34 方法）、契约微扩 2 处（D3 记账）、**零迁移零权限种子**。

**Tests**: 宪法 IV 走 TDD（先红后绿；测试与被测包同目录 `promotion_impl_test.go` / `activity_impl_test.go`）；C 端联动用既有 C 端方法断言。

## Phase 1: Foundational（微扩、fixture）

- [ ] T001 `promotion.go` 接口微扩：`ICouponLogic` +`AdminDetail`（D3-① 记账）；`model.CouponTemplate` +`ValidType`（D3-② 记账）
- [ ] T002 [P] 测试基座 `promotion_impl_test.go`：券/满减 fixture 与自清（精确名、子表先删、含 C 端联动断言所用数据）

## Phase 2: US1 券模板管理（P1）🎯 MVP

**独立验收**: 创建→列表/详情→停发→C 端消失→记录下钻；类型/有效期条件校验。

- [ ] T003 [P] [US1] 测试（红）：创建矩阵（type/validType 条件、金额非法、start≥end）/列表含已领数/详情/AdminUpdate 停发 → C 端首页可领券联动/软删（user_coupon 不受影响）/记录分页
- [ ] T004 [US1] 实现 `promotion_impl.go` 的 `AdminList/AdminDetail/AdminCreate/AdminUpdate/AdminDelete/AdminRecords`；连线 6 桩（券 6 端点清零）
- [ ] T005 [US1] 实现期把批次 09 记账的"满减 scope 下单计价"保持挂账（不动 order_impl）

## Phase 3: US2 满减活动（P1）

- [ ] T006 [P] [US2] 测试（红）：档位+范围嵌套创建往返/重复门槛 50008/范围三型校验（全场无 target、2/3 必须 target）/全量替换/C 端列表 scopeDesc 联动/时间窗
- [ ] T007 [US2] 实现 `activity_impl.go` 的 `FullReduction{List,Create,Detail,Update,Delete}`；连线 5 桩

## Phase 4: US3 拼团/秒杀/砍价（P1）

- [ ] T008 [P] [US3] 测试（红）：三类活动 CRUD + SetItems 全量替换/唯一键撞 1062 转业务码/秒杀已售保护（订单项引用+sold_count>0 拒绝移除）/砍价 original>floor/拼团 groupPrice 必填/C 端列表联动
- [ ] T009 [US3] 实现 `activity_impl.go` 的 `GroupBuy*/FlashSale*/Bargain*`（含死契约 Detail×3, D4）；连线 15 桩

## Phase 5: US4 助力（P2）

- [ ] T010 [P] [US4] 测试（红）：rewardType 两型落库（券→reward_ref 校验存在；积分→config JSON）/perLimit>1 接受且如实回显/requiredCount≥1/C 端 rewardDesc 联动
- [ ] T011 [US4] 实现 `Assist{List,Create,Update,Delete}`；连线 4 桩

## Phase 6: Polish & 批次收尾

- [ ] T012 N+1 收敛自查（D7）：列表档位/范围/SPU 简介批量 IN 查询；`internal/routes/` 或 controller 层补**端点级可达性**测试（admin 权限挂载生效：无权限 403 类/有权限可达——沿 014/015 wiring 先例）
- [ ] T013 `go test ./...` 两连跑全绿（01~09 零退化）；`golangci-lint run` **0 issues**
- [ ] T014 `check-stub` 对账：**admin 58→28**（四渠道合计 69→39）；quickstart 冒烟 5 组
- [ ] T015 更新 `specs/PROGRESS.md` 批次 10 状态 ✅ 与完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001→T002→US1→US2→US3→US4→收尾（同文件域, 顺序执行; 测试任务可先行红）
- 满减（US2）先于 US3 无硬依赖, 但同文件按序避免冲突

## Notes（本批硬约束）

- 写入防线铁律（D1）与全量替换已售保护（D2）不可妥协
- 契约微扩两处（T001）**先记账后动码**（PROGRESS §五 D3 行随开工行已记）
- 权限挂载按路由注释（FR-6）；权限点已就位勿新增迁移
- 禁止手改生成物；不实现拼团成团/退款、发奖实际发放、下单侧满减认范围（挂账批次 13/横切）
- C 端联动断言复用既有 C 端方法（`IMarketingLogic`/`publicCouponBriefs` 口径）
