---
description: "任务清单：评价域（批次 08）"
---

# Tasks: 评价域（014-review）

**Feature**: `specs/014-review` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/review-endpoint-mapping.md](./contracts/review-endpoint-mapping.md)

**Input**: 4 端点（shop，全为 `CodeNotImplemented` 桩）、5 个方法（4 复用 + 1 微扩 `ProductList`）、**零迁移、零新 DTO、零生成物变更**。

**Tests**: 宪法 IV 走 TDD（`service` 全部新写方法先红后绿；测试与被测包同目录）。

## Phase 1: Foundational（契约微扩与基座）

- [x] T001 `internal/service/shop/review.go` 接口微扩 `ProductList`（返回 汇总 + 分页列表；见 contracts 的签名），其余三方法签名不动；注释按 D1/D5 口径更新（V1 自动通过 = 通过后展示）
- [x] T002 `internal/service/shop/review_impl.go` 骨架：`ReviewLogicImpl` + `NewReviewLogic()`（形态跟随 shop 域既有 struct 实现）
- [x] T003 `internal/service/shop/review_impl_test.go` 测试基座：可复用的"已完成订单 + 订单项（含 spu_name/sku_specs 快照）"构造与评价 seed/cleanup（复用 `aftersale_impl_test.go` 的 fixture 构造方式；清理精确匹配、子表先删）

## Phase 2: US1 提交评价（P1）🎯 MVP

**独立验收**：已完成订单项可评（返回 ID 且 audit_status=1）；重复/并发提交只成功一次；他人/未完成/越界被拒。

- [x] T004 [P] [US1] 测试（红）：`Create` —— 成功落库（快照 spu_name/sku_specs、images 数组、audit_status=1）、重复提交 40010 且活库仅 1 条、8 并发恰好 1 次成功（`uk_order_item` 兜底 + 1062 转业务码）、非本人按不存在、订单未完成 **40006**（状态不允许）、评分越界 10001
- [x] T005 [US1] 实现 `Create`（归属/已完成/未评过校验 → 快照落库；捕获唯一键冲突 1062 → 40010）
- [x] T006 [US1] 连线端点 `internal/controller/shop/shop_v1_review_create.go`（桩清零 ×1；`requireMember` + `parseID`）

## Phase 3: US2 商品评价列表与汇总（P1）

**独立验收**：只出审核通过；星级筛选生效；汇总与列表同源同筛；匿名/非匿名脱敏正确。

- [x] T007 [P] [US2] 测试（红）：`ProductList` —— 2 通过 + 1 待审 + 1 驳回 → 只出 2 条；`score` 筛选；汇总 Total/Avg/分布只含通过；新提交（自动通过）后 Total+1；匿名显示"匿名用户"、非匿名脱敏、空昵称占位；规格/回复随条目返回
- [x] T008 [US2] 实现 `ProductList`（汇总 + 列表同源同筛；脱敏 helper 在 shop 域内——`maskNickname` 属兄弟域不可 import，注释说明；昵称读 `user` 表，跟随 `order_mgmt_impl.go` 既有先例）
- [x] T009 [US2] 连线端点 `internal/controller/shop/shop_v1_product_review_list.go`（桩清零 ×1）

## Phase 4: US3 我的评价（P2）

**独立验收**：只返回本人（含全部审核状态）；追评/回复一并透出。

- [x] T010 [P] [US3] 测试（红）：`MyList` —— 本人 2 条（不同状态）+ 他人 1 条 → 只出本人的 2 条且状态正确；含追评/回复的条目一并返回；分页生效
- [x] T011 [US3] 实现 `MyList`
- [x] T012 [US3] 连线端点 `internal/controller/shop/shop_v1_my_review_list.go`（桩清零 ×1）

## Phase 5: US4 追加评价（P2）

**独立验收**：30 天前主评可追评一次；二次/超期/他人/空内容一律拒绝。

- [x] T013 [P] [US4] 测试（红）：`Extra` —— 成功（extra_content/extra_time 落库）、二次 40011、超 90 天 40011、他人按不存在、空内容 10001、并发追评只成功一次（条件更新）
- [x] T014 [US4] 实现 `Extra`（条件更新：`extra_content=''` + 90 天内；affected=0 时回查区分"已追评/超期"）
- [x] T015 [US4] 连线端点 `internal/controller/shop/shop_v1_review_extra.go`（桩清零 ×1）

## Phase 6: Polish & 批次收尾

- [x] T016 `go test ./...` 两连跑全绿（含 01~07 既有链路零退化）；`golangci-lint run` 保持 **0 issues**
- [x] T017 `make check-stub` 对账：**shop 16→12**（本批 4 端点全清）；按 [quickstart.md](./quickstart.md) 冒烟 8 项
- [x] T018 更新 `specs/PROGRESS.md` 批次 08 状态 ✅ 与本批完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001 → T002 → T003 → 各故事
- US1（T004~T006）→ US2（T007~T009）→ US3（T010~T012）→ US4（T013~T015）
- US2/US3 都要"评价人展示"helper（US2 先建，US3 复用）
- 收尾依赖全部故事完成

## Notes（本批硬约束）

- **写入防线**：一项一评靠 `uk_order_item` 唯一索引兜底，冲突 1062 **必须转业务码**（40010）而非裸抛；追评用**条件更新**（`extra_content=''`）——批次 07 的 C1/C2 教训（读-改窗口）不得复发
- **口径一致**：商品评价列表与汇总**同源同筛**（只出 audit_status=1）；我的评价不筛审核状态
- **V1 自动通过**（用户裁定 B）：新评价 `audit_status=1`；**合规债务已记账**（公开上线前必须补后台审核面）
- 批次文件边界：`internal/service/shop/review*.go`、4 个 shop controller 桩、本批 specs 目录、`specs/PROGRESS.md`
- 禁止手改生成物；**不动**支付/订单/售后域既有实现；**不实现** `Reply`/`Audit`/`PendingAudit`（死契约）

## 完成记录（2026-09-22）

- **4 端点全清**: `check-stub` 等价对账 **shop 16→12**（四渠道合计 85→81）；SC-001 达标
- **验证**: `go test ./...` 连续两次全绿（01~07 既有链路零退化）；`golangci-lint run` **0 issues**；**零迁移**（表与索引既有）
- **新增 9 个测试函数**（`review_impl_test.go`）: 提交（快照/自动通过/重复/他人/未完成/越界）、**并发提交只成功一次**（唯一键兜底 + 1062 转业务码）、商品列表（只出通过/星级筛选/汇总同源/匿名脱敏/规格快照）、我的评价（含全部状态/追评/回复/分页/不含他人）、追评（成功/二次/超期 90 天边界/他人/空内容/**并发只成功一次**）
- **契约微扩（记账）**: `IReviewLogic` 补 `ProductList`（第 4 个端点在 api 存在而接口漏定义；同批次 01/02/04 先例）
- **口径**: V1 自动通过（用户裁定 B，见 §五 开工行）；展示侧严格只出 audit_status=1
- **新发现（记账待裁定）**: 追评请求的 `images` 参数**既无存储列（表无 extra_images）也无出参字段**（`ReviewItem`/`MyReviewItem` 均无 extraImages）→ 契约自相矛盾；本批**不为其新增列**（那会写入永远无法展示的数据），该字段当前被忽略，已在实现处注释说明
- **未实现（死契约，记账）**: `Reply`/`Audit`/`PendingAudit` 三方法无 API 端点
