---
description: "Task list for 商品目录域实现"
---

# Tasks: 商品目录域实现（product-catalog）

**Input**: Design documents from `/specs/005-product-catalog/`（spec/plan/research/data-model/quickstart）+ 既有 `internal/service/shop` 接口定义（IProductLogic/IInventoryLogic）

**Prerequisites**: plan.md (required), spec.md (required), research.md（D1 冗余列/D2 dao 扩表/D3 口径/D4 权限过渡/D5 搜索/D6 测试策略/D7 足迹）

**Tests**: TDD 红绿为强制约定（宪法 IV + §五）——每故事先写失败测试再实现；测试配置确定性注入（gdb/gredis SetConfig），数据自建+前置清理（research D6）。

**Organization**: 执行顺序按数据流依赖调整——后台管理（US2，数据入口）先于浏览（US1，数据消费）；US3 库存与 US4 搜索随后。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件、无依赖）
- **[Story]**: 归属用户故事
- service 实现落位 `internal/service/shop/product_impl.go` / `inventory_impl.go`（`_impl` 后缀与接口定义文件区分）

---

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 编写迁移 `migrations/000034_spu_price_range.up.sql`：product_spu 增 `price_min/price_max DECIMAL(10,2) NULL` + `idx_price_min`（research D1/data-model §一）+ down 占位
- [ ] T002 扩表生成：`gf gen dao -t "product_category,product_brand,product_spu,product_sku,inventory,inventory_log,product_review"`（research D2），确认 entity/do/dao 七表生成物无手改
- [ ] T003 [P] 种子超管过渡（research D4）：bootstrap 增 EnsureSuperAdmin——启动时 `admin_user` 无任何账号则按 env ADMIN_USERNAME/ADMIN_PASSWORD 创建 `is_super=1` 账号（bcrypt），env 缺省则跳过并输出警告日志

**Checkpoint**: 冗余列就绪、七表 dao 生成、超管过渡形态明确

---

## Phase 2: User Story 2 - 后台商品管理 (Priority: P1) 🎯 数据入口

**Goal**: 分类/品牌/SPU/SKU 全套管理能力（数据入口，先于浏览）

**Independent Test**: 管理接口序列"建分类→建品牌→建 SPU→加 2 SKU→上架"成功；30004/30005/30006 三防御正反断言；规格组合并发冲突映射 30006

- [ ] T001~T003 见上；管理接口权限说明：后台登录/权限点真实校验属治理特性，本特性管理接口以 Auth 白名单外直通形态运行（operator 固定 `admin:seed`），**治理特性接入时补校验**——记入本文件 Notes
- [ ] T004 [US2] 先写失败测试 `internal/service/shop/product_impl_test.go`（CategoryTree/CRUD/引用禁删 30004、品牌 CRUD/名称唯一、SPU 上下架 30005、规格组合 30006、限售设置全量替换语义）
- [ ] T005 [US2] 实现 `internal/service/shop/product_impl.go`：IProductLogic 全部方法（entity→model DTO 转换、规格组合唯一捕获 1062→30006、上架校验 30005、分类引用禁删 30004、限售 JSON 全量替换）；Where/Data 一律 do 对象
- [ ] T006 [US2] SKU 变更价格冗余刷新：`refreshSpuPriceRange(spuId)`（同事务重算启用 SKU min/max），接入 SKU 创建/编辑/启停/删除四个路径
- [ ] T007 [US2] 控制器接线：`internal/controller/admin` 商品/分类/品牌相关桩填充→调 service（api/admin/v1 既有契约），写接口挂权限点常量（product:*）
- [ ] T008 [US2] 转绿：`make build && go test ./internal/service/shop/` 全绿 + 三防御正反用例通过

**Checkpoint**: 商品数据入口可用；三防御有测试证据

---

## Phase 3: User Story 3 - 库存查询与调整 (Priority: P2)

**Goal**: 库存闸门能力（列表/调整留痕/预警）

**Independent Test**: 调整 -20 拒绝 30007；+50 后列表/预警同步且 log +1（operator 可追溯）

- [ ] T009 [P] [US3] 先写失败测试 `internal/service/shop/inventory_impl_test.go`（列表可售推导/调整正负/负向超调 30007/留痕条数/operator/预警过滤）
- [ ] T010 [US3] 实现 `internal/service/shop/inventory_impl.go`：IInventoryLogic 全部方法（调整=条件更新+同事务 log change_type=5；预警=可售≤warn_count）
- [ ] T011 [US3] 控制器接线：admin 库存三桩填充 + 转绿

**Checkpoint**: 库存闸门可用（ADR-0001 留痕口径落地）

---

## Phase 4: User Story 1 - 商品浏览 (Priority: P1)

**Goal**: 游客/会员浏览闭环（依赖 US2 数据）

**Independent Test**: quickstart 场景一全序列 + 下架反向（30001 语义）

- [ ] T012 [P] [US1] 先写失败测试（internal/service/shop）：分类树过滤禁用/软删、列表筛选排序（五档）与价格区间闭区间、详情可售口径（SPU上架 AND SKU启用 AND 可售>0）、评价汇总（audit_status=1 聚合）、足迹 UPSERT（FR-006）
- [ ] T013 [US1] 实现浏览方法（product_impl.go 补 C 端侧）+ 价格区间空值形态（全部 SKU 禁用 → 卡片价格区间空 + 不可下单标注）
- [ ] T014 [US1] 控制器接线：common/shop 渠道浏览桩填充（categories/brands/products/{id}/search/reviews）+ 转绿

**Checkpoint**: 游客"逛"契约闭环（SC-004 价格/排序断言过）

---

## Phase 5: User Story 4 - 限售区域与搜索 (Priority: P3)

**Goal**: 限售合规 + 搜索体验

**Independent Test**: 限售标记在详情可见；关键词命中+销量排序正确；特殊字符转义

- [ ] T015 [P] [US4] 搜索实现：关键词 LIKE（%/_/' 转义）+ 筛选排序复用（research D5），`product_impl.go` 补 Search
- [ ] T016 [US4] 测试：关键词命中/排序正确/特殊字符不报错（转义验证）；限售标记详情断言（spec US4 场景 1）

**Checkpoint**: 搜索与限售能力齐备

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T017 quickstart 四场景终验（curl 序列）输出留证到本目录 verification.md
- [ ] T018 宪法 IV 终验：`make build && make test` 全绿；`gofmt -l api internal` 空输出
- [ ] T019 管理组权限点真实校验登记 backlog（依赖后台登录特性）；治理特性接入时补 operator 真实身份

---

## Dependencies & Execution Order

- Setup（T001~T003）→ **US2（T004~T008，数据入口）** → US3（T009~T011）→ US1（T012~T014，消费数据）→ US4（T015~T016）→ Polish
- 说明：spec 优先级 US1 在前，但执行顺序按数据流调整（浏览验收需要管理端产生的数据）；两者优先级同为 P1
- [P] 任务可并行；T005/T006 同文件严格串行

### Parallel Opportunities

- T003 与 T001/T002 并行；T009 测试编写与 US2 实现并行；US3/US4 测试先行与实现解耦

---

## Implementation Strategy

### MVP First

Setup → US2 全部 = 后台可建商品数据（第一可演示增量：管理面可用）。

### Incremental Delivery

US2 → US3（库存闸门）→ US1（浏览开放）→ US4（合规+搜索）：每阶段一次中文提交、测试随行、可独立评审。

---

## Notes

- 分层契约 §二：Where/Data 一律 do 对象、返回 entity→转 model DTO、禁 g.Map（实现阶段违反即打回）
- 管理接口鉴权为占位形态（Auth 直通）——operator 记 `admin:seed`，治理特性接入时统一补真实校验（本文件 Notes 已登记）
- specs_hash 为 DB 生成列：插入自动计算，无需应用层处理；并发冲突捕获 1062 → 30006
- 验证以输出为证；每任务组一次提交
