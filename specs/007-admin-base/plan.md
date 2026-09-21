# Implementation Plan: 后台账户与系统配置（007-admin-base）

**Branch**: `007-admin-base` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

实现平台治理域首批：admin 渠道登录会话、后台账号/角色/权限管理、系统配置管理面 + common 基础端点，共 22 端点
（admin 20 + common 2）的 service 实现（`internal/service/system` 包级函数）与 controller 连线。
关键技术点：会话增加 audience 维度（修复跨渠道越权，research D1）、controller 显式权限点挂接（D2）、
种子超管迁移 000035（D3）。零表结构变更、零新第三方依赖。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（gdb + dao 生成物 + do 强类型）

**Primary Dependencies**: 既有栈复用——`library/security`（SessionManager，本批加 audience）、
`library/captcha`（登录验证码校验）、`library/sms`（mock:sms:{phone} 读取）、`internal/errcode`；零新增

**Storage**: MySQL 8.4（admin_user/admin_role/admin_permission/admin_user_role/admin_role_permission/
system_config/admin_login_log 七表全既有）；Redis（会话 audience key、mock 短信）
新迁移：`000035_admin_seed`（纯数据种子：超管账号 1 个，research D3）

**Testing**: `internal/service/system/*_test.go` gtest 红绿（005/006 基座模式：确定性配置注入 + 数据自建清理）；
`auth_flow_test.go` 适配 audience 签名后全量回归

**Target Platform**: `internal/service/system`（RBAC 13 + Auth 2+1Profile + Config 2）、
`internal/middleware`（Auth 渠道判定扩展 + RequirePerm helper）、controller 四文件组、`migrations/000035`

**Project Type**: GoFrame 模块化单体·治理域纵切片

**Performance Goals**: 管理端低频操作，无专项指标（HasPermission 每请求联查可接受，缓存为演进项）

**Constraints**: 分层契约 §二/§四（do/entity 强类型、controller 只调 service/middleware）；
TDD 红绿；宪法 V（不新建顶层包——RequirePerm 落既有 middleware 包，种子落既有 migrations 目录）

**Scale/Scope**: 22 端点、service 包级函数约 20 个、DTO 仅新增 AdminProfile 1 个、测试覆盖全部行为级验收场景

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位 service/system 既有包与既有 controller 渠道；middleware 横切层调 service 包级函数（无反向依赖） |
| II 统一技术栈 | ✅ | 零新依赖；bcrypt 复用既有加密库（service/user 先例） |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | 红绿 TDD + quickstart 序列 + `make check-stub` 数字验收 |
| V 简单优先 | ✅ | 不建顶层包/抽象层；权限挂接选显式 helper（D2 否决路由映射表）；读取组件推迟（D5） |
| 工程约束 §迁移 | ✅ | 000035 为纯数据种子（000032 先例），down 可回退，migrate-fresh 可重放 |
| 注销权/合规 | ✅ | 软删+禁用语义既有；账号删除为软删留痕（审计要求） |

**Phase 1 复查**：✅ 四产物（research/data-model/contracts/quickstart）无越界；唯一 spec 假设修订（无迁移→
加 000035 数据种子）已按 PROGRESS 防偏离条款 4 在变更记录记账。

## Project Structure

### Documentation (this feature)

```text
specs/007-admin-base/
├── plan.md                        # 本文件
├── research.md                    # Phase 0：D1~D8 决策与勘察
├── data-model.md                  # Phase 1：实体/会话/权限数据流
├── contracts/admin-permissions.md # Phase 1：22 端点 × 权限点映射 + 错误码
├── quickstart.md                  # Phase 1：验证指南
├── checklists/requirements.md     # specify 阶段质量单
└── tasks.md                       # Phase 2（/speckit-tasks 生成，本批后续）
```

### Source Code (repository root)

```text
apps/server/
├── migrations/
│   ├── 000035_admin_seed.up.sql     # 种子超管（新增, 纯数据）
│   └── 000035_admin_seed.down.sql
├── internal/
│   ├── library/security/session.go  # SessionManager 加 audience（D1, 含既有调用点适配）
│   ├── middleware/
│   │   ├── auth.go                  # 渠道判定（/admin→admin会话）+ 管理员账号态校验
│   │   └── permission.go            # RequirePerm helper（新增, D2）
│   ├── service/system/
│   │   ├── rbac.go / ops.go         # 接口（契约, 微扩 Profile）
│   │   ├── rbac_impl.go             # 角色/权限/账号管理实现（新增）
│   │   ├── auth_impl.go             # 登录/改密/Profile 实现（新增）
│   │   ├── config_impl.go           # 配置 List/Update 实现（新增）
│   │   └── *_impl_test.go           # TDD 红绿（新增）
│   ├── model/dto_system.go          # 补 AdminProfile（唯一 DTO 增量）
│   ├── controller/admin/            # auth/rbac/config 三组 20 桩填充 + common 2 桩
│   └── errcode/                     # 增补权限不足等新码（见 contracts）
└── tests 回归: service/user/auth_flow_test.go（audience 适配）
```

**Structure Decision**: 全部落位既有分层与既有包（宪法 V）；唯一新文件为
`middleware/permission.go`（横切 helper 既有包内）、三个 service/system 实现文件、000035 迁移对。

## Complexity Tracking

> 无违例——不填。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
