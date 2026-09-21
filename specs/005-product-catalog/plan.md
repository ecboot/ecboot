# Implementation Plan: 商品目录域实现（product-catalog）

**Branch**: `005-product-catalog` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

实现商品目录域：C 端浏览（分类树/品牌/列表/详情/搜索/评价汇总/足迹软触）+ 后台商品管理（分类/品牌/SPU/SKU 全套 CRUD 与两级上下架/限售）+ 库存（查询/调整/预警）。落位 `internal/service/shop`（IProductLogic/IInventoryLogic 实现）+ `internal/controller` 接线 + 新迁移 000034（SPU 价格冗余列）。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（gdb 数据访问、dao 生成物复用）

**Primary Dependencies**: 既有栈；`gf gen dao` 扩表生成（product_category/brand/spu/sku/inventory/inventory_log/product_review 七表）

**Storage**: MySQL 8.4；新增迁移 000034：`product_spu` 增 `price_min/price_max DECIMAL(10,2) NULL` 冗余列（列表价格区间与价格排序的索引友好的实现，SKU 变更时同事务刷新）+ 两列索引调整

**Testing**: service 层 gtest（确定性配置注入，红绿 TDD）；数据策略=测试内建数据+清理（不依赖种子）；端到端 curl 序列（quickstart）

**Target Platform**: `internal/service/shop`（IProductLogic/IInventoryLogic 实现）、`internal/controller/{common,shop,admin}` 接线、`internal/routes` 权限点挂接

**Project Type**: 业务域实现（契约→service→仓储→测试全链）

**Performance Goals**: 商品列表/详情 p95 < 200ms（本地容器；列表走冗余列+idx 无聚合 JOIN）

**Constraints**: 分层契约 §二（Where/Data 用 do、返回 entity→转 model DTO、禁 g.Map 传参）；service 仅调 dao+manager；TDD 红绿；管理写接口挂权限点

**Scale/Scope**: 七表 dao 生成、1 迁移、2 个 logic 实现（约 30 方法）、3 渠道控制器接线、service 层测试全覆盖（SC-006）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 实现 004 已定义的 IProductLogic/IInventoryLogic，零新顶层包 |
| II 统一技术栈 | ✅ | gdb + gf gen dao；搜索 V1 为 DB 条件查询（ES 演进不改契约） |
| III 中文优先 | ✅ | 注释/文档中文 |
| IV 可验证交付 | ✅ | 红绿 TDD + `make build/test` + quickstart 序列以输出为证 |
| V 简单优先 | ✅ | 价格冗余列而非查询时聚合（索引友好）；搜索 V1 LIKE 不引 ES |
| 工程约束 §二 强类型 | ✅ | Where/Data 一律 do 对象；返回 entity→转 model DTO；禁 g.Map |
| 工程约束 §五 TDD | ✅ | 先测试后实现（service 层全覆盖 SC-006） |
| B2C+多门店 | ✅ | 单店数据模型（store 域不涉本特性） |

**Phase 1 复查**：✅ research/data-model/contracts/quickstart 无越界；000034 为本特性唯一迁移。

## Project Structure

### Documentation (this feature)

```text
specs/005-product-catalog/{plan,research,data-model,quickstart}.md, tasks.md(后续)
```

### Source Code (apps/server)

```text
migrations/000034_spu_price_range.up/down.sql   # price_min/price_max 冗余列
internal/service/shop/
├── product_impl.go      # IProductLogic 实现（浏览+后台管理; 转换 entity↔model DTO）
├── inventory_impl.go    # IInventoryLogic 实现（含 ADR-0001 留痕）
└── product_impl_test.go # 红绿测试（数据内建+清理）
internal/controller/{common,shop,admin}   # 桩填充→调 service
internal/routes/route.go                  # 权限点挂接（admin 组写接口）
internal/consts/permission.go             # 权限点常量核对（已就绪）
```

**Structure Decision**: 实现文件以 `_impl` 后缀与接口定义文件（product.go 接口聚合）区分；仓储直接使用 gf 生成 dao（不建 repository 中间层——dao 已是强类型仓储，评审约定 service 仅调 dao+manager）。

## Complexity Tracking

> 无宪法违例。
