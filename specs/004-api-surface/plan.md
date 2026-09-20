# Implementation Plan: 全角色 API 接口面设计（api-surface）

**Branch**: `004-api-surface` | **Date**: 2026-09-20 | **Spec**: [spec.md](./spec.md)

## Summary

依据 68 张业务表，为游客/会员/管理员三角色设计四渠道完整 API 契约（约 140 个端点），落入 `apps/server/api` 的 GoFrame 接口定义层（v1 Req/Res 结构体 + IUser 接口聚合 + 路由注册），控制器与业务实现为后续特性。契约六要素（路径/方法/鉴权/入参/出参/错误码）逐端点完整。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（接口定义层：`g.Meta` 路由元 + `v:"..."` 校验规则 + `dc:` 字段说明）

**Primary Dependencies**: 既有工程结构（api/{渠道}/v1 三元组约定，跟随用户已生成的代码形态）；零新第三方

**Storage**: 无直接存储访问（契约层）；错误码/权限点/统一响应沿用 003 契约与 `internal/errcode` 既有实现

**Testing**: `make build`（契约层编译零错误）+ 路由自检（启动后路由表含全部注册路径）+ 契约走查（quickstart 三旅程核对）

**Target Platform**: `apps/server/api/{common,user,shop,admin}`（现有目录四渠道）

**Project Type**: 接口契约层设计（GoFrame api 定义）

**Performance Goals**: N/A（契约层）

**Constraints**: 渠道目录互不引用；统一响应由中间件包装（Res 结构体只定义 data 部分）；B2B2C 零接口（宪法）；003 已设计端点引用不重做但纳入总路由

**Scale/Scope**: 约 140 端点契约、四渠道 × (v1 定义 + 渠道接口聚合文件)；零迁移、零业务实现

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 契约落位既有四渠道目录；渠道互不引用由目录约定保证 |
| II 统一技术栈 | ✅ | GoFrame 标准接口定义范式（与用户已生成代码同构），零新依赖 |
| III 中文优先 | ✅ | summary/dc 全中文 |
| IV 可验证交付 | ✅ | `make build` + 路由自检 + 契约走查三证 |
| V 简单优先 | ✅ | 契约复用既有 errcode/统一响应；不预设未实现能力的结构（如导出仅留扩展声明） |
| 工程约束 | ✅ | B2C+多门店：零 B2B2C 接口；store 接口为线下载体语义 |
| 开发流程 | ✅ | 003 已设计 8 端点引用不重述；RefreshToken 等用户已落地的实现语义反向纳管契约 |

**Phase 1 复查**：✅ contracts 四渠道文件与结构树一致；无越界（无业务逻辑、无迁移）。

## Project Structure

### Documentation (this feature)

```text
specs/004-api-surface/
├── plan.md / research.md / data-model.md / quickstart.md
└── contracts/
    ├── conventions.md      契约总则（响应/分页/鉴权/错误码/命名）
    ├── common-api.md       公共渠道端点表
    ├── user-api.md         会员渠道端点表
    ├── shop-api.md         商城渠道端点表
    └── admin-api.md        管理渠道端点表
```

### Source Code (apps/server/api)

```text
api/
├── common/            common.go(ICommonV1 聚合) + v1/{captcha,store,ping}.go
├── user/              user.go(IUserV1) + v1/{auth,profile,address,favorite,
│                      footprint,point,coupon,message,distribution,share}.go
├── shop/              shop.go(IShopV1) + v1/{product,activity,cart,order,pay,
│                      aftersale,review,bargain,assist,banner}.go
└── admin/             admin.go(IAdminV1) + v1/{auth,product,category,brand,
                       inventory,order,aftersale,coupon,promotion,banner,store,
                       logistics,member,distribution,risk,rbac,audit,config,dashboard}.go
路由注册: internal/routes（分组挂载 + 鉴权中间件标注）
```

**Structure Decision**: 一域一 v1 文件；Req/Res 成对定义；鉴权级别由路由分组（公开组/会员组/管理组）承载而非结构体标注；权限点以常量登记于 `internal/consts` 供 RBAC 权限树种子使用。

## Complexity Tracking

> 无宪法违例。边界项：管理端导出接口仅留扩展声明（FR 假设）。
