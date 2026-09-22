---
description: "任务清单：营销后台（批次 10）"
---

# Tasks: 营销后台（016-marketing-admin）

**Feature**: `specs/016-marketing-admin` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/marketing-admin-endpoint-mapping.md](./contracts/marketing-admin-endpoint-mapping.md)

**Input**: 30 端点（admin，全为桩）、`ICouponLogic`+`IActivityLogic` 补实现（34 方法）、契约微扩 2 处（D3 记账）、**零迁移零权限种子**。

**Tests**: 宪法 IV 走 TDD（先红后绿；测试与被测包同目录 `promotion_impl_test.go` / `activity_impl_test.go`）；C 端联动用既有 C 端方法断言。

## Phase 1: Foundational（微扩、fixture）

- [x] T001 `promotion.go` 接口微扩：`ICouponLogic` +`AdminDetail`（D3-① 记账）；`model.CouponTemplate` +`ValidType`（D3-② 记账）
- [x] T002 [P] 测试基座 `promotion_impl_test.go`：券/满减 fixture 与自清（精确名、子表先删、含 C 端联动断言所用数据）

## Phase 2: US1 券模板管理（P1）🎯 MVP

**独立验收**: 创建→列表/详情→停发→C 端消失→记录下钻；类型/有效期条件校验。

- [x] T003 [P] [US1] 测试（红）：创建矩阵（type/validType 条件、金额非法、start≥end）/列表含已领数/详情/AdminUpdate 停发 → C 端首页可领券联动/软删（user_coupon 不受影响）/记录分页
- [x] T004 [US1] 实现 `promotion_impl.go` 的 `AdminList/AdminDetail/AdminCreate/AdminUpdate/AdminDelete/AdminRecords`；连线 6 桩（券 6 端点清零）
- [x] T005 [US1] 实现期把批次 09 记账的"满减 scope 下单计价"保持挂账（不动 order_impl）

## Phase 3: US2 满减活动（P1）

- [x] T006 [P] [US2] 测试（红）：档位+范围嵌套创建往返/重复门槛 50008/范围三型校验（全场无 target、2/3 必须 target）/全量替换/C 端列表 scopeDesc 联动/时间窗
- [x] T007 [US2] 实现 `activity_impl.go` 的 `FullReduction{List,Create,Detail,Update,Delete}`；连线 5 桩

## Phase 4: US3 拼团/秒杀/砍价（P1）

- [x] T008 [P] [US3] 测试（红）：三类活动 CRUD + SetItems 全量替换/唯一键撞 1062 转业务码/秒杀已售保护（订单项引用+sold_count>0 拒绝移除）/砍价 original>floor/拼团 groupPrice 必填/C 端列表联动
- [x] T009 [US3] 实现 `activity_impl.go` 的 `GroupBuy*/FlashSale*/Bargain*`（含死契约 Detail×3, D4）；连线 15 桩

## Phase 5: US4 助力（P2）

- [x] T010 [P] [US4] 测试（红）：rewardType 两型落库（券→reward_ref 校验存在；积分→config JSON）/perLimit>1 接受且如实回显/requiredCount≥1/C 端 rewardDesc 联动
- [x] T011 [US4] 实现 `Assist{List,Create,Update,Delete}`；连线 4 桩

## Phase 6: Polish & 批次收尾

- [x] T012 N+1 收敛自查（D7）：列表档位/范围/SPU 简介批量 IN 查询；`internal/routes/` 或 controller 层补**端点级可达性**测试（admin 权限挂载生效：无权限 403 类/有权限可达——沿 014/015 wiring 先例）
- [x] T013 `go test ./...` 两连跑全绿（01~09 零退化）；`golangci-lint run` **0 issues**
- [x] T014 `check-stub` 对账：**admin 58→28**（四渠道合计 69→39）；quickstart 冒烟 5 组
- [x] T015 更新 `specs/PROGRESS.md` 批次 10 状态 ✅ 与完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001→T002→US1→US2→US3→US4→收尾（同文件域, 顺序执行; 测试任务可先行红）
- 满减（US2）先于 US3 无硬依赖, 但同文件按序避免冲突

## Notes（本批硬约束）

- 写入防线铁律（D1）与全量替换已售保护（D2）不可妥协
- 契约微扩两处（T001）**先记账后动码**（PROGRESS §五 D3 行随开工行已记）
- 权限挂载按路由注释（FR-6）；权限点已就位勿新增迁移
- 禁止手改生成物；不实现拼团成团/退款、发奖实际发放、下单侧满减认范围（挂账见 PROGRESS §六 P1——13 批收官后归口独立横切批）
- C 端联动断言复用既有 C 端方法（`IMarketingLogic`/`publicCouponBriefs` 口径）

## 完成记录（2026-09-22）

- **30 端点全清**: `check-stub` 对账 **admin 58→28**, 四渠道合计 **69→39**
- **实现**: 新写 `promotion_impl.go`（券 7 方法含微扩 AdminDetail）+ `activity_impl.go`（五类 28 方法）; 连线 30 桩; **契约微扩 D3-①~⑧**（接口 +AdminDetail; DTO +ValidType/Status×5/列表玩法字段×5——api 有而 DTO 漏的最小适配, 全部记账）
- **实现期新发现（测试当场抓到）**: ①拼团表时间列为 `valid_start_at/valid_end_at`（与其余活动表 `start_time/end_time` 不同名）→ 通用函数时间列参数化; ②**同值 UPDATE affected=0 误判"不存在"**（MySQL 不计同值行）→ 存在性前置 Count + 更新幂等成功（批次 02 修改语义先例）
- **防线落地**: 满减档位 50008（同门槛唯一, 建库+1062 双路）; 场次商品 uk_activity_sku 1062→50002; 秒杀**已售保护**（事务内 LockUpdate 复核 "sold_count>0 的行不可移除"）; 全量替换事务先删后插; 金额 money.FromYuanString 拒 >2 位小数
- **端点级测试**: `marketing_admin_wiring_test.go`——30 端点全量探活（超管直通不回 10003/10005）+ 未登录对照必须 10003（RequirePerm 挂接漏洞的守卫）
- **C 端联动（SC-3）**: 停发/软删 → C 端首页可领券/满减列表/秒杀/助力列表即时消失（复用 C 端既有过滤口径, 断言钉住）
- **验证**: `go test ./...` 两连跑全绿; `golangci-lint` **0 issues**; 零迁移零权限种子（000032 已就位, 勘察核对）
