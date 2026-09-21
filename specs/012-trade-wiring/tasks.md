# Tasks: 交易闭环连线（012-trade-wiring）

**Input**: Design documents from `/specs/012-trade-wiring/`
**Prerequisites**: plan.md、spec.md、research.md（D1~D5）、data-model.md、contracts/trade-endpoint-mapping.md、quickstart.md
**Tests**: 宪法 IV——全部新写方法 TDD 红绿；连线以"可调用 + 既有行为回归"验收。

## Phase 1: Setup

- [ ] T001 基线确认：`go build ./...` 与 `go test ./...` 全绿（批次 05 收口态）；确认桩数 shop 35 / admin 69 / user 13

## Phase 2: 前置修复（阻塞 US2，FR-001）

- [ ] T002 [P] `internal/service/shop/promotion_calc.go`：`calcPointDeductFen` 移除 `AND deleted=0`（该表无此列，批次 05 评审债务）；补行为级断言（有余额勾选→>0 / 未勾选→0 / 无账户行→0 不报错）于既有 trade_test.go 或新文件
- [ ] T003 `go test ./internal/service/shop/...` 绿（既有交易链路零退化）

## Phase 3: US2 购物车 + US3 订单 C 端（连线 + Confirm 新写）

- [ ] T004 [US3] `OrderLogicImpl.Confirm`（30→40 条件更新 + 日志 + 佣金事件位；非 30 → 40006）+ TDD
- [ ] T005 [US2] 连线购物车 5 端点（桩清零 ×5；Checkout 的积分抵扣由 T002 修复承接）
- [ ] T006 [US3] 连线订单 C 端 4 端点（create/list/detail/cancel）+ Confirm 端点（桩清零 ×5）
- [ ] T007 `go test ./...` 绿

## Phase 4: US4 支付（新写，P1）

- [ ] T008 [P] [US4] `internal/service/shop/pay_impl_test.go`（红）：Create（应付<=0 拒绝）/Status；
      **HandlePayNotify 幂等四层**（条件更新 affected=0 → 幂等应答；金额不符 → 拒绝+留档；订单已取消 → 不推进）；
      成功推进（支付单20/订单20/库存核销/余额消费完成）；HandleRefundNotify 幂等推进售后单
- [ ] T009 [US4] `internal/service/shop/pay_impl.go`（新增）：5 方法（含 CloseExpired：过期支付单关闭）
- [ ] T010 [US4] 连线 shop 支付 4 端点（桩清零 ×4；notify 两端点公开——白名单已有）
- [ ] T011 `go test ./...` 绿

## Phase 5: US5 后台订单（新写，P1）

- [ ] T012 [P] [US5] `internal/service/shop/order_mgmt_test.go`（红）：AdminList（状态/关键词/分页）/
      AdminDetail/Deliver（仅 20 可发 + **停用物流公司拒绝** + 物流信息落库）/AdminCancel（同 C 端语义）/
      SellerRemark（**C 端 Detail 不含该列**）/CancelTimeout/AutoConfirm
- [ ] T013 [US5] `order_impl.go` 追加 8 方法（Confirm 已在 T004）
- [ ] T014 [US5] 连线 admin 5 端点（桩清零 ×5；deliver 挂 order:deliver、cancel 挂 order:cancel、remark 挂 order:update）
- [ ] T015 `go test ./...` 绿

## Phase 6: US6 会员优惠券（新写，P2）

- [ ] T016 [P] [US6] `internal/service/user/coupon_impl_test.go`（红）：AvailableTemplates（过滤领完/超限/停发）/
      Receive（**同事务防超发+限领**；超限 → 50001）/Mine（四态 + **惰性过期判定**）/
      UsableForOrder（门槛 + 抵扣降序）/Consume（绑单号）/ReturnBack（退回，有效期不变）
- [ ] T017 [US6] `internal/service/user/coupon_impl.go`（新增，包级函数 6 个）
- [ ] T018 [US6] 连线 user 3 端点（桩清零 ×3）
- [ ] T019 `go test ./...` 绿

## Phase 7: Polish & 批次收尾

- [ ] T020 `make test` 全绿 + golangci-lint 本批文件零问题（含既有交易链路零退化——SC-004）
- [ ] T021 `make check-stub` 对账：shop 35→21、admin 69→64、user 13→10；按 quickstart 冒烟；
      更新 PROGRESS 批次 06 状态 ✅ 与完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001 → T002/T003（前置修复, 阻塞结算试算）→ US3（Confirm）→ US2/US3 连线 → US4 支付 → US5 后台 → US6 券 → 收尾
- 各故事文件边界：pay_impl.go / order_impl.go / coupon_impl.go 互不重叠
- 收尾依赖全部故事完成

## Notes

- 禁止手改 `internal/dao`、`internal/model/entity|do`（本批零迁移，无生成物变更）
- **支付回调幂等四层防线为硬约束**（资金安全）；金额不符必留档
- 本批文件边界：plan.md「Source Code」+ specs/012-trade-wiring/ + specs/PROGRESS.md
