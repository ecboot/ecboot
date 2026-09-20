---
description: "Task list for 全角色 API 接口面设计（契约层落码）"
---

# Tasks: 全角色 API 接口面设计（契约层落码）

**Input**: Design documents from `/specs/004-api-surface/`（plan、spec、research、data-model、contracts/ 四渠道细表、quickstart）

**Prerequisites**: plan.md (required), spec.md (required), contracts/ (逐端点契约权威)

**Tests**: 契约层验证 = `make build` 编译 + 路由自检 + quickstart 三旅程走查（每故事收尾任务），无业务测试（实现属后续特性）。

**Organization**: 按用户故事分组（US1 游客 / US2 会员交易 / US3 会员资产社交 / US4 后台运营 / US5 后台治理）。一域一任务（一个 v1 文件 + 控制器桩 + 路由注册）；不同渠道文件可并行。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件、无依赖）
- **[Story]**: 归属用户故事
- 每个域任务隐含三件事：① 按 contracts 对应文件定义 v1 Req/Res（含 g.Meta 路由元/校验规则/dc 中文说明）；② `gf gen ctrl` 生成控制器桩；③ 在 `internal/routes` 对应分组注册路由。收尾任务统一编译验证。

---

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 修订 003 既有端点路径至新分组前缀：`api/common/v1/captcha.go`（/common/captcha*）、`api/user/v1/auth.go`（/user/login/sms、/user/login/wx，Res 增补 isNew/userId 字段）；同步 `internal/routes` 分组挂载，删除根路径旧注册
- [x] T002 [P] 在 `internal/errcode` 登记 004 错误码段位常量（3xxxx 商品 / 4xxxx 交易 / 5xxxx 促销 / 6xxxx 分销 / 7xxxx 门店 / 8xxxx 后台治理，逐语义命名，见 data-model 登记表）
- [x] T003 [P] 在 `internal/consts` 登记权限点常量全集（模块:资源:动作，约 60 个，清单见 data-model §三），作为 admin_permission 种子与路由权限校验共用
- [x] T004 建立 api 层公共契约结构：`api/common/v1/common.go` 增补分页嵌入结构（PageReq/PageRes，ID/金额 string 约定注释）；各渠道 v1 文件引用

**Checkpoint**: 基础约定就绪（错误码/权限点/分页/路径前缀），域契约可开始

---

## Phase 2: Foundational (Blocking Prerequisites)

- [x] T005 建立四渠道路由分组骨架：`internal/routes` 挂载 `/common`（公开）、`/user`（公开子组 login + 会员子组 Bearer 中间件）、`/shop`（公开子组 + 会员子组）、`/admin`（公开子组 login + 管理子组 Bearer+权限点中间件）；中间件未实现处留接口占位（编译通过）

**Checkpoint**: 分组就绪——所有域路由有归属组

---

## Phase 3: User Story 1 - 游客浏览接口 (Priority: P1) 🎯 MVP

**Goal**: 游客可逛：商品/活动/运营位/门店公开契约

**Independent Test**: `make build` 通过 + 路由表含全部公开端点 + 游客旅程路由可达（quickstart 验证三旅程 1）

- [x] T006 [P] [US1] 定义 `api/common/v1/store.go`：门店列表（区县/经纬度+距离排序）与详情（contracts/common-api.md 门店节）
- [x] T007 [P] [US1] 定义 `api/common/v1/ping.go` + `api/common/v1/share.go`：探针与分享上报（可选会员身份）
- [x] T008 [P] [US1] 定义 `api/shop/v1/product.go`：分类树/品牌/商品列表/详情（含 SKU+评价汇总）/搜索/评价列表（contracts/shop-api.md 商品浏览节）
- [x] T009 [P] [US1] 定义 `api/shop/v1/activity.go`：拼团/秒杀/砍价/助力/满减五类活动浏览（contracts/shop-api.md 营销活动节）
- [x] T010 [P] [US1] 定义 `api/shop/v1/banner.go`：轮播与楼层公开接口
- [x] T011 [US1] 生成控制器桩（gf gen ctrl：common/store|ping|share + shop/product|activity|banner）、路由注册至公开组、`make build` 验证、更新渠道聚合接口文件

**Checkpoint**: 游客旅程契约完备——"逛"的接口面闭环

---

## Phase 4: User Story 2 - 会员交易接口 (Priority: P1)

**Goal**: 会员能买能退：地址/购物车/订单/支付/售后/评价契约

**Independent Test**: 路由表含全部交易端点；会员组路由无凭证返回 10003（中间件占位时以挂载存在为证）

- [x] T012 [P] [US2] 定义 `api/user/v1/address.go`：地址 CRUD + 设默认（含三级区划码字段）
- [x] T013 [P] [US2] 定义 `api/shop/v1/cart.go`：购物车五接口（含结算试算）
- [x] T014 [US2] 定义 `api/shop/v1/order.go`：创建（四玩法上下文+幂等凭证+余额抵扣）/列表/详情/取消/确认收货（contracts/shop-api.md 订单节）
- [x] T015 [US2] 定义 `api/shop/v1/pay.go`：发起支付/状态查询/支付回调/退款回调（回调=开放组+签名验证语义标注）
- [x] T016 [P] [US2] 定义 `api/shop/v1/aftersale.go`：申请/列表/详情/撤销/寄回单号
- [x] T017 [P] [US2] 定义 `api/shop/v1/review.go`：提交/追评/我的评价
- [x] T018 [US2] 生成控制器桩、路由注册（会员组）、`make build` 验证

**Checkpoint**: 交易闭环契约完备（"能买能退"）

---

## Phase 5: User Story 3 - 会员资产与社交接口 (Priority: P2)

**Goal**: 会员资产/社交/分销契约

**Independent Test**: 路由表含 /user 全部资产端点；分销端点齐全（对照 contracts/user-api.md 分销节 11 行）

- [x] T019 [P] [US3] 定义 `api/user/v1/profile.go`：资料/等级成长值/积分账户与流水
- [x] T020 [P] [US3] 定义 `api/user/v1/favorite.go` + `api/user/v1/footprint.go`：收藏（含复活语义）/足迹
- [x] T021 [P] [US3] 定义 `api/user/v1/coupon.go`：可领列表/领券/我的券
- [x] T022 [P] [US3] 定义 `api/user/v1/message.go`：站内信/已读/偏好
- [x] T023 [US3] 定义 `api/user/v1/distribution.go`：分销中心 11 端点（申请/状态/关系/比例/记录/账户/流水/提现/邀请/推广码）
- [x] T024 [US3] 生成控制器桩、路由注册、`make build` 验证

**Checkpoint**: 会员资产与社交契约完备

---

## Phase 6: User Story 4 - 后台运营管理接口 (Priority: P2)

**Goal**: 后台运营全域契约（每写接口带权限点常量引用）

**Independent Test**: 路由表含 /admin 运营端点全集；权限点与 T003 常量一一对应

- [x] T025 [P] [US4] 定义 `api/admin/v1/auth.go`：后台登录/refresh/登出/个人资料/改密
- [x] T026 [P] [US4] 定义 `api/admin/v1/product.go` + `api/admin/v1/category.go` + `api/admin/v1/brand.go`：商品/分类/品牌管理（含 SKU 与限售区域）
- [x] T027 [P] [US4] 定义 `api/admin/v1/inventory.go`：库存列表/调整/预警
- [x] T028 [P] [US4] 定义 `api/admin/v1/order.go`：订单管理五接口（含卖家备注）
- [x] T029 [P] [US4] 定义 `api/admin/v1/aftersale.go`：售后管理六接口
- [x] T030 [US4] 定义 `api/admin/v1/promotion.go`：券/满减/拼团/秒杀/砍价/助力管理（本特性最大单文件，按 contracts/admin-api.md 促销节逐行核对）
- [x] T031 [P] [US4] 定义 `api/admin/v1/operation.go` + `api/admin/v1/store.go` + `api/admin/v1/logistics.go`：运营位/门店/物流字典
- [x] T032 [P] [US4] 定义 `api/admin/v1/distribution.go` + `api/admin/v1/member.go`：分销管理/会员管理
- [x] T033 [US4] 生成控制器桩、路由注册（管理组+权限点引用）、`make build` 验证

**Checkpoint**: 后台运营契约完备

---

## Phase 7: User Story 5 - 后台治理接口 (Priority: P3)

**Goal**: RBAC/审计/风控/配置/看板契约

**Independent Test**: 路由表含治理端点；权限点 system:*/risk:*/dashboard:read 齐备

- [x] T034 [P] [US5] 定义 `api/admin/v1/rbac.go`：角色/权限树/分配/账号管理
- [x] T035 [P] [US5] 定义 `api/admin/v1/audit.go` + `api/admin/v1/risk.go`：审计查询/风控规则与事件/申诉
- [x] T036 [P] [US5] 定义 `api/admin/v1/config.go` + `api/admin/v1/dashboard.go`：系统配置/三看板
- [x] T037 [US5] 生成控制器桩、路由注册、`make build` 验证

**Checkpoint**: 全角色契约完备

---

## Phase 8: Polish & Cross-Cutting Concerns

- [x] T038 路由自检：启动应用核对路由总数与 contracts 端点行数一致（含 mock 条件注册语义）；将端点统计回填本文件 Notes
- [x] T039 RBAC 种子：由 T003 权限点常量生成 `admin_permission` 种子 SQL（`migrations/000033_rbac_seed.up.sql`），与权限树接口数据同源
- [x] T040 契约走查终验：执行 quickstart 三旅程 + 一致性抽检；`make build && make test` 全绿留证

---

## Dependencies & Execution Order

- **Setup（T001~T004）→ Foundational（T005）→ US1~US5（T006~T037）→ Polish（T038~T040）**
- 域契约任务全部 [P] 可并行（不同文件）；各故事的"收尾任务"（T011/T018/T024/T033/T037）串行（编译验证 + 路由注册冲突检查）
- T014 依赖 T013（订单创建语义引用购物车试算字段）；T030 依赖 T003（权限点常量）

### Parallel Opportunities

- T002/T003/T004 并行；US1~US5 的域任务在 T005 完成后全量并行（10+ 文件互不依赖）
- 渠道间无依赖：common/user/shop/admin 可由不同执行者并行推进

---

## Implementation Strategy

### MVP First (US1)

T001~T005 → T006~T011：游客浏览契约闭环即可演示"接口面"价值（路由表+OpenAPI 文档随 GoFrame 自动生成）。

### Incremental Delivery

每故事 = 一批 v1 文件 + 控制器桩 + 路由注册 + 编译验证 + 一次提交（如 `feat: 商城渠道游客浏览接口契约`）。US1→US2→US3→US4→US5 按优先级交付，各批独立可评审。

---

## Notes

- 契约六要素权威 = contracts/ 四文件；v1 结构体字段以 dc 中文说明承载语义，实现阶段（controller/service）不得私自变更契约
- 003 已有端点（T001）只做路径迁移与字段增补，行为逻辑不动（业务实现在 003 特性范围）
- 每任务或逻辑组完成后提交一次（中文 Conventional Commits）
- 验证以输出为证（make build / 路由表 / quickstart）
