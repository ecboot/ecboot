---
description: "任务清单：营销 C 端（批次 09）"
---

# Tasks: 营销 C 端（015-marketing-c）

**Feature**: `specs/015-marketing-c` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/marketing-c-endpoint-mapping.md](./contracts/marketing-c-endpoint-mapping.md)

**Input**: 12 端点（shop，全为桩；**本批收官 shop 渠道**）、3 个新接口/实现、~10 个新 DTO、2 个新端口、1 处跨批改动（秒杀欠账清偿）、**零迁移**。

**Tests**: 宪法 IV 走 TDD（先红后绿；测试与被测包同目录）；另补**端点级可达性测试**（014 评审 C1 的教训）。

## Phase 1: Foundational（契约、DTO、端口、基座）

- [x] T001 `internal/service/shop/marketing.go` / `bargain.go` / `assist.go`：新建三个接口（签名见 contracts §一），注释写明口径（只出有效活动/按结束时间升序/一人一刀/一助力/不砍穿/端口降级）
- [x] T002 `internal/model/dto_shop.go`：新增 C 端 DTO（五类公开项 + 场次商品摘要 + 砍价/助力进度视图与结果 + `IndexAggregate`；见 data-model §四）
- [x] T003 `ports.go`：新增 `IRiskHit`（风控，未装配降级放行 + 告警）与 `IAssistReward`（发奖意图，未装配告警降级）
- [x] T004 测试基座（`marketing_impl_test.go`/`play_impl_test.go`）：造"五类活动 + 场次商品 + 会员 + 商品库存"的可复用 fixture 与自清（精确匹配 + 子表先删 + 清理 helper/inventory_log）

## Phase 2: US1 秒杀（P1，含批次 07 欠账清偿）🎯 MVP

**独立验收**：列表可见（含秒杀价/剩余/限购）→ 秒杀价下单（快照正确 + 双库存扣减）→ 支付回调推进且核销成功 → 取消双回补。

- [x] T005 [P] [US1] 测试（红）：`PublicFlashSales`（进行中+预告、排序、分页）+ 秒杀下单链路（秒杀价快照、活动/商品双库存条件更新、场次 SKU 不符 → 50002、活动库存不足 → 40003、商品库存不足 → 40001）+ 取消回补 `sold_count`/`locked` + **支付回调核销成功**（批次 07 欠账的直接验证）
- [x] T006 [US1] 实现 `marketing_impl.go` 的 `PublicFlashSales`（与首页入口共用查询）
- [x] T007 [US1] **跨批清偿**：`order_impl.go` —— `collectLines` 取 `flash_price` 快照（`price`=秒杀价/`original_price`=原价）、步骤 3 增加活动库存条件更新 + `(activity_id,sku_id)` 校验、**删除批次 07 的秒杀拦截闸**；`cancelBy` 回补 `sold_count`
- [x] T008 [US1] 连线 `shop_v1_flash_sale_list.go`（桩清零 ×1）

## Phase 3: US2 拼团列表（P1）

- [x] T009 [P] [US2] 测试（红）：`PublicGroupBuys` 只出进行中、含成团价与人数字段、按结束时间升序、分页
- [x] T010 [US2] 实现 `PublicGroupBuys` + 连线 `shop_v1_group_buy_list.go`（桩清零 ×1）

## Phase 4: US3 砍价（P1）

**独立验收**：发起（首刀）→ 他人帮砍 → 到底价 → 状态可下单；一人一刀/本人不可/不砍穿/超时/并发。

- [x] T011 [P] [US3] 测试（红）：`PublicBargains`（列表）+ `Launch`（首刀金额 = 均分向下取整、发起次数上限 50006）+ `Progress`（进度与帮砍列表脱敏）+ `Cut`（一人一刀 50004、本人不可 50002、**末刀补齐到恰好底价**、超时 50002、**并发只成功一次**、风控端口拦截 10005）
- [x] T012 [US3] 实现 `bargain_impl.go`（三方法）
- [x] T013 [US3] 连线 4 端点：`shop_v1_bargain_activity_list.go` / `shop_v1_bargain_launch.go` / `shop_v1_bargain_progress.go` / `shop_v1_bargain_cut.go`（桩清零 ×4）

## Phase 5: US4 助力（P2）

**独立验收**：发起 → 一人一助力 → 达标发奖意图仅一次；发起次数上限/本人不可/并发。

- [x] T014 [P] [US4] 测试（红）：`PublicAssists`（列表）+ `Launch`（发起次数上限 50006、活动时间窗 50002）+ `Progress`（进度与助力人列表脱敏）+ `Help`（一人一助力 50007、本人不可 50002、**达标投递发奖意图且仅一次**、并发达标只投递一次、风控端口拦截 10005）
- [x] T015 [US4] 实现 `assist_impl.go`（三方法）
- [x] T016 [US4] 连线 4 端点：`shop_v1_assist_list.go` / `shop_v1_assist_launch.go` / `shop_v1_assist_progress.go` / `shop_v1_assist_help.go`（桩清零 ×4）

## Phase 6: US5 满减列表 + US6 首页聚合（P2）

- [x] T017 [P] [US5] 测试（红）：`PublicFullReductions` 只出进行中、含档位与范围摘要
- [x] T018 [US5] 实现 `PublicFullReductions` + 连线 `shop_v1_full_reduction_list.go`（桩清零 ×1）
- [x] T019 [P] [US6] 测试（红）：`Index` 一次返回轮播+楼层+五入口（各 ≤3 条）+券；**任一块为空返回空数组**；`api/shop/v1/index.go` 契约扩写
- [x] T020 [US6] 实现 `Index`（复用既有 `PublicBanners`/`PublicFloors`/`PublicList` 与本批列表）+ 连线 `shop_v1_index.go`（桩清零 ×1）

## Phase 7: Polish & 批次收尾

- [x] T021 `internal/routes/` 补**端点级可达性测试**（公开端点游客可达 + 会员端点合法凭证可达/无凭证 10003；`/shop/index` 的白名单归属按探活决定）
- [x] T022 `go test ./...` 两连跑全绿（含 01~08 既有链路零退化，**特别注意交易链路**）；`golangci-lint run` 保持 **0 issues**
- [x] T023 `make check-stub` 对账：**shop 12→0**（shop 渠道收官）；按 [quickstart.md](./quickstart.md) 冒烟 5 组
- [x] T024 更新 `specs/PROGRESS.md` 批次 09 状态 ✅ 与本批完成 commit（同 commit）并提交

## Dependencies & Execution Order

- T001（接口）→ T002（DTO）→ T003（端口）→ T004（基座）→ 各故事
- US1（T005~T008）→ US2（T009~T010）→ US3（T011~T013）→ US4（T014~T016）→ US5（T017~T018）→ US6（T019~T020）
- US6 的首页入口依赖各列表方法（US1~US5 先行）
- 收尾依赖全部故事完成

## Notes（本批硬约束）

- **秒杀双库存**（批次 07 教训）: 活动库存（`flash_sale_item` 条件更新）与商品库存（`inventory` 锁）**都要做**，任一不足即拒绝；取消回补两者；**支付回调不得因核销未命中而回滚**（这是本批清偿的直接验收点）
- **写入防线**: 一人一刀/一人一助力靠 `uk_record_helper` 唯一键兜底 + 1062 转业务码；计数与价格推进用条件更新 + 判 `RowsAffected`（012/013/014 三轮评审确立，不得复发）
- **口径一致**: 首页入口与各列表**共用同一查询函数**（批次 08 的"列表与汇总漂移"教训）
- **端口降级**: 风控/发奖端口未装配时放行/告警，不阻塞玩法可用
- 批次文件边界：`service/shop/{marketing,bargain,assist}*.go`、`ports.go`、`internal/model/dto_shop.go`、`api/shop/v1/index.go`、12 个 shop controller 桩、`internal/routes/`（可达性测试）、本批 specs 目录、`specs/PROGRESS.md`
- **跨批改动（已记账）**: `order_impl.go` 的秒杀分支与 `cancelBy`（批次 07 ledger 明确移交）
- 禁止手改生成物；**不实现**拼团成团/解散/退款、风控规则引擎、发奖实际发放（端口留出口）

## 完成记录（2026-09-22）

- **12 端点全清**: `check-stub` 对账 **shop 12→0（渠道收官）**；四渠道合计 81→69；SC-001 达标
- **验证**: `go test ./...` 两连跑全绿（01~08 零退化）；`golangci-lint` **0 issues**；新增迁移 000040 已应用（版本 40）
- **秒杀欠账清偿（批次 07 移交）**: 秒杀价快照（`price`=秒杀价/`original_price`=原价）、活动+商品**双库存**条件更新、取消双回补、**删除批次 07 的临时拦截**；`TestFlashSalePayCallback` 直接钉住"支付回调核销成功"（批次 07 的失败点）
- **实现期新发现（测试当场抓到）**: ①白名单前缀吞掉会员动作（`/{id}/cut`、`/{id}/helpers`）→ 新增 `publicGetOnlyPrefixes`（同前缀只放行 GET）+ 端点级可达性测试（实证能抓回归）；②`/shop/index` 漏配白名单；③助力达标判定用旧计数（并发下永不达标）→ 自增后重读
- **契约微扩（记账）**: `IndexRes` 占位→聚合结构；`BargainSku` +itemId/maxCutCount；`ActivitySkuBrief` +ItemId
- **未做（边界）**: 拼团成团/解散/退款（批次 10/11）；发奖实际发放（端口留出口）；风控规则引擎（批次 12）；砍价发起次数限制（表无字段）
