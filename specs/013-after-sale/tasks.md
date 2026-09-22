---
description: "任务清单：售后域（批次 07）"
---

# Tasks: 售后域（013-after-sale）

**Feature**: `specs/013-after-sale` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/after-sale-endpoint-mapping.md](./contracts/after-sale-endpoint-mapping.md)

**Input**: 11 端点（shop 5 + admin 6，全为 `CodeNotImplemented` 桩）、11 个新写方法、**唯一迁移 000038**（两列）、0 新表、0 新 DTO。

**Tests**: 本批按宪法 IV 走 TDD（`service` 全部新写方法先红后绿；测试与被测包同目录）。

## Phase 1: Setup（迁移与生成物）

- [x] T001 新增 `migrations/000038_after_sale_operator.{up,down}.sql`：给 `after_sale_order` 追加 `operator_id VARCHAR(64) NOT NULL DEFAULT ''` 与 `fail_reason VARCHAR(255) NOT NULL DEFAULT ''`（依据 research D3）；应用到应用库并核对版本=38
- [x] T002 `gf gen dao` 重生成 `after_sale_order` 生成物（`internal/dao/internal/after_sale_order.go`、`internal/model/{entity,do}/after_sale_order.go`）；**保留真变更、还原换行符噪音**（同 012 修订轮做法）

## Phase 2: Foundational（基座与口径）

- [x] T003 `internal/service/shop/aftersale_impl_test.go` 测试基座：可复用的"已完成订单 + 订单项 + 库存"构造与售后单 seed/cleanup（复用 `order_create_test.go` 的地址/商品 fixture；清理一律精确匹配 + 子表先删）
- [x] T004 `internal/service/shop/aftersale.go` 接口注释按裁定修正：撤销边界收窄为"待审核/待寄回/待退款"（research D4、spec FR-013），并在注释中点明"退款中不可撤"的理由
- [x] T005 `internal/service/shop/aftersale_impl.go` 骨架：`AfterSaleLogicImpl` + `NewAfterSaleLogic()`（形态跟随 shop 域既有 struct 实现）

## Phase 3: US1 会员申请售后（P1）🎯 MVP

**独立验收**：已完成订单的订单项可申请（返回单号+状态 10）；超额/超量/他人项被拒且不落库。

- [x] T006 [P] [US1] `aftersale_impl_test.go`（红）：`Apply` 校验矩阵——本人已完成订单项成功（type=1/2、refund_amount 按行实付折算到分）、数量超行数量→40008、累计超额→40008（含"已拒绝/已撤销不占用额度"）、他人订单项→40009、非可售后状态→40008
- [x] T007 [US1] `aftersale_impl.go` 实现 `Apply`（含 salesNo 生成、可退数量累计口径、退款金额计算与"末笔取剩余全额"消尾差）
- [x] T008 [US1] 连线 shop 申请端点 `internal/controller/shop/shop_v1_after_sale_create.go`（桩清零 ×1；userId 经会员守卫）

## Phase 4: US2 后台审核（P1）

**独立验收**：仅退款同意→发起退款（40）；退货退款同意→待寄回（20）；拒绝→90 且原因必填；非 10 →40006。

- [x] T009 [P] [US2] 测试（红）：`Approve` 两分支（type=1 直达 40、type=2 到 20）、`Reject`（原因必填、落 reject_reason）、状态机守卫（非 10 →40006）、`operator_id` 落库
- [x] T010 [US2] 实现 `AdminList` / `AdminDetail` / `Approve` / `Reject` + 退款发起内部方法（`paychannel.Channel.Refund`，幂等键=售后单号；成功→40，失败→停 30 写 `fail_reason`）
- [x] T011 [US2] 连线 admin 4 端点（`admin_v1_admin_after_sale_{list,detail,approve,reject}.go`；权限点按 api 注释 `aftersale:audit`，**须先核对 RBAC 是否可装载该权限点**，缺失会 10005）

## Phase 5: US3 退货寄回与确认收货（P1）

**独立验收**：待寄回可填单号（仅 type=2 且 20）；未填单号确认收货→拒绝；确认后进入 40。

- [x] T012 [P] [US3] 测试（红）：`SubmitReturn`（仅 20 且 type=2；非 20 →40006）、`ConfirmReceipt`（缺单号→40008；有单号→30→40）
- [x] T013 [US3] 实现 `SubmitReturn` / `ConfirmReceipt`
- [x] T014 [US3] 连线 shop 寄回端点 `shop_v1_after_sale_logistics.go` + admin 确认收货 `admin_v1_admin_after_sale_confirm_receipt.go`（桩清零 ×2）

## Phase 6: US4 退款闭环与重试（P1）

**独立验收**：40 收到成功回调→50；重复回调不改写；未匹配回调被拒且留档；失败停在 30 可重试；50 不可再退；0 元不调渠道直收口。

- [x] T015 [P] [US4] 测试（红）：回调推进（复用既有 `HandleRefundNotify`）→50 + refund_time；重复回调幂等；`RetryRefund`（30→40、清 fail_reason、50→40006 拒绝）；0 元边界（不调渠道、直接 50、refund_no 空）
- [x] T016 [US4] 实现 `RetryRefund` 与 0 元收口分支（含 `fail_reason` 清理）
- [x] T017 [US4] 连线 admin 退款重试 `admin_v1_admin_after_sale_retry_refund.go`（桩清零 ×1；权限 `aftersale:refund`）

## Phase 7: US5 会员查询与撤销（P2）

**独立验收**：列表按状态筛选+分页；仅本人可见；撤销边界 {10,20,30} 可撤且额度释放，40/终态→40006。

- [x] T018 [P] [US5] 测试（红）：`List`（状态筛选+分页+仅本人）、`Detail`（他人→40009）、`Cancel`（10/20/30 可撤→91 且可退数量恢复；40→40006；50→40006）
- [x] T019 [US5] 实现 `List` / `Detail` / `Cancel`
- [x] T020 [US5] 连线 shop 3 端点（`shop_v1_after_sale_{list,detail,cancel}.go`；桩清零 ×3）

## Phase 8: US6 完成后的账实一致（P2）

**独立验收**：退货退款完成→库存回补 quantity、仅退款不回补；订单 refund_status 按"完成态口径"重算；完成时投递佣金冲销事件。

- [x] T021 [P] [US6] 测试（红）：库存回补（仅 type=2；`total` 增加、`locked` 不变）、`refund_status` 三段（0/1/2）、佣金冲销事件出口被调用
- [x] T022 [US6] 实现完成副作用（回补库存 / 重算 `refund_status` / 投递冲销事件位）
- [x] T023 `go test ./...` 全绿（含批次 01~06 既有链路零退化）

## Phase 9: Polish & 批次收尾

- [x] T024 `golangci-lint run` 保持 **0 issues**；`make test` 全绿
- [x] T025 `make check-stub` 对账：**shop 21→16、admin 64→58**（本批 11 端点全清）；按 [quickstart.md](./quickstart.md) 冒烟 7 项
- [ ] T026 更新 `specs/PROGRESS.md` 批次 07 状态 ✅ 与本批完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001 → T002 → T003（基座需要新列）→ T005（骨架）→ 各故事
- US1（T006~T008）→ US2（T009~T011）→ US3（T012~T014）→ US4（T015~T017）→ US5（T018~T020）→ US6（T021~T023）
- US2/US3 都要"发起退款"内部方法：US2 先建，US3 复用（不重复实现）
- US4 依赖 US2/US3 已能进入 30/40；US6 依赖 US2/US3 的完成路径存在
- 收尾依赖全部故事完成

## Parallel Opportunities

- T006 / T009 / T012 / T015 / T018 / T021（各故事的测试文件段落可按故事并行编写），但 `aftersale_impl.go` 为同一文件 → 实现任务串行

## Notes（本批硬约束）

- **状态迁移一律"条件更新 + 判 RowsAffected"**（affected=0 → 40006）：这是批次 06 修正后的既有惯例，I1/I11 同型缺陷不得复发
- 金额一律 `money` 库 + 分运算；`refund_amount ≤ 订单项 pay_amount`；0 元不调渠道
- 退款幂等键 = 售后单号（`out_refund_no = after_sale_no`）；回调推进复用既有 `HandleRefundNotify`（**不改它**）
- 撤销边界：{10,20,30} 可撤，"退款中"不可撤（用户裁定 D4）
- 批次文件边界：`migrations/000038_*`、`internal/service/shop/aftersale*.go`、11 个 controller 桩、本批 specs 目录、`specs/PROGRESS.md`
- 禁止手改生成物（`internal/dao`、`internal/model/entity|do`）；禁止改 `refund_notify` 与支付/订单域既有实现

## 完成记录（2026-09-22）

- **11 端点全清**: `make check-stub` 等价对账 **shop 21→16、admin 64→58**（四渠道合计 96→85）；SC-001 达标
- **验证**: `go test ./...` 连续两次全绿（批次 01~06 既有链路零退化）；`golangci-lint run` **0 issues**；迁移全量重放（`ecboot_fresh`）到版本 38 且 `after_sale_order.operator_id` 就位
- **新增测试 21 个函数**（`aftersale_impl_test.go`）: 申请校验矩阵/额度释放/尾差消解/审核两分支/状态机守卫/寄回与确认收货/退款回调幂等/渠道失败留痕与重试/0 元边界/撤销边界/列表详情归属/完成副作用（回补·退款状态·佣金冲销事件）
- **实现期偏离（已记账）**:
  1. 新增**迁移 000038**（原 spec 假设零迁移）——`IAfterSaleLogic` 四个后台方法都收 `operator` 而表无列可承载（research D3）
  2. **跨文件钩子**: `pay_impl.go` 的 `HandleRefundNotify` 改为事务包裹并在条件更新命中时调用 `afterSaleFinishedSideEffects`——不改则**渠道回调完成时永不回补库存/不更新订单退款状态**（SC-005 无法成立）；tasks 原写"不改它"，此为必要修正
  3. `ports.go` 新增 `ICommissionReverse` 事件出口（结算属批次 11，本批只投递）
  4. DTO 微扩: `AfterSaleSummary`/`AfterSaleDetail` 各 +`UserId`（后台列表/详情契约要求，C 端不映射；同批次 04"按契约最小适配"先例）
  5. 权限点**无需新迁移**——`aftersale:read/audit/refund` 已在 000032 种子与 `consts/permission.go` 中就位（实现前核对，未凭猜）
- **口径落地**: 撤销边界 {10,20,30}（用户裁定）；退款发起时机=进入待退款的同一动作；0 元不调渠道直收口
