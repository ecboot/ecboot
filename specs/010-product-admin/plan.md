# Implementation Plan: 商品目录后台与 C 端浏览收口（010-product-admin）

**Branch**: `010-product-admin` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

把 005 已交付的商品域服务实现接到契约端点（24 个：admin 19 + shop 5），新写库存后台（3 个），
并按权限点挂接管理端 22 个端点。**不改既有业务逻辑**（005/006 已验收语义零回退）。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3

**Primary Dependencies**: 既有栈复用；零新依赖

**Storage**: MySQL 8.4（product_category/brand/spu/sku/inventory 全既有）；零新迁移

**Testing**: `internal/service/shop/inventory_impl_test.go`（库存三方法 TDD 红绿，复用同包基座）；
连线端点的验收以"可调用 + 权限 + 既有行为"，由既有 005 测试（product_impl_test.go）承担语义回归

**Target Platform**: `internal/service/shop`（inventory_impl 新增）、controller
`admin_v1_admin_{category,brand,spu,sku,inventory}_*.go` ×22 + `shop_v1_{category,brand,product,search}*.go` ×5

**Project Type**: GoFrame 模块化单体·商品域连线收口

**Performance Goals**: 无专项（既有实现已交付）

**Constraints**: 分层契约（controller 只调 service）；**调用形态跟随 005 遗产**——controller 经
`shop.NewProductLogic()` 调 struct 方法（与同域既有形态一致, 不引入包级函数双轨）；
TDD 仅施于新写部分（库存三方法）；连线部分以"端到端可调用 + 权限"验收

**Scale/Scope**: 27 端点、新写 3 方法 + 1 测试文件、22+5 个 controller 桩填充

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落既有 service/shop 域与既有 controller 渠道；无跨域/跨渠道引用 |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | 库存新写 TDD；连线以 `check-stub` 数字 + 端到端冒烟 + 既有测试回归为证 |
| V 简单优先 | ✅ | 不重写既有实现；调用形态跟随既有 struct（不引入双轨）；零迁移零新 DTO（既有 DTO 76 个已覆盖） |
| 工程约束 §二/五 | ✅ | controller 只做绑定+组装；service 用 do/entity 强类型 |
| 跨域协作 | ✅ | 零跨域 |

**Phase 1 复查**：✅ 映射表（contracts）覆盖 27 端点无一遗漏；新写范围严格限定库存三方法。

## Project Structure

### Documentation (this feature)

```text
specs/010-product-admin/
├── plan.md / research.md / data-model.md / quickstart.md
├── contracts/product-endpoint-mapping.md   # 27 端点 × 既有方法 × 权限码（连线清单）
├── checklists/requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── internal/
│   ├── service/shop/
│   │   ├── inventory_impl.go           # 新增: List/Warnings/Adjust（含流水留痕）
│   │   └── inventory_impl_test.go      # 新增: TDD 红绿
│   ├── controller/admin/               # category ×4 + brand ×4 + spu ×7 + sku ×4 + inventory ×3 桩填充
│   ├── controller/shop/                # category_tree + brand_list + product_list/detail/search 桩填充
│   └── （既有 product_impl.go / browse_impl.go 不改业务逻辑）
└── specs/PROGRESS.md
```

**Structure Decision**: 唯一新文件 `inventory_impl.go` + 测试；其余为桩填充（27 个 controller）。

## Complexity Tracking

> 无违例。连线工作量集中在 controller 层映射，但每条映射都有既有 impl 与契约双向对照。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
