# Implementation Plan: 物流公司与运营装修（009-logistics-ops）

**Branch**: `009-logistics-ops` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

实现物流公司字典与运营装修（banner/楼层）共 15 端点（admin 13 + shop 2）的 service 实现与 controller 连线。
关键技术点：装修域接口与 DTO 新建（D1，非微扩）、楼层 config 约定与商品摘要装配（D2/D3，失效剔除）、
投放时段 SQL 侧判定（D4）。零表结构变更、零新迁移、1 个新错误码（40012）、DTO 新增 7 个。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（gdb + dao 生成物 + do 强类型）

**Primary Dependencies**: 既有栈复用（零新依赖）；同包 helper 复用（`firstImage`）

**Storage**: MySQL 8.4（logistics_company 000015、operation_banner/operation_floor 000026 全既有）；零新迁移

**Testing**: `internal/service/shop/operation_impl_test.go` gtest 红绿（复用该包既有 init 基座与清理惯例）

**Target Platform**: `internal/service/shop`（logistics_impl/operation_impl）、controller
`admin_v1_admin_{logistics,banner,floor}_*.go` ×13 + `shop_v1_{banner,floor}_list.go` ×2

**Project Type**: GoFrame 模块化单体·交易支撑与内容运营纵切片

**Performance Goals**: 装修内容读多写少（C 端高频读）；在投过滤走 `idx_position_status` 与 `idx_status_sort`

**Constraints**: 分层契约 §二/§四（do/entity 强类型、controller 禁触 dao）；TDD 红绿；宪法 V（不新建域）

**Scale/Scope**: 15 端点、service 包级函数 13 个 + 2 个微扩/新建接口、DTO 7 个、错误码 1 个

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位 service/shop 既有域与既有 controller 渠道；跨渠道零引用；商品装配同域直读（不跨域事件） |
| II 统一技术栈 | ✅ | 零新依赖（JSON 处理用标准库 encoding/json，与 browse_impl 同源） |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | TDD 红绿 + quickstart 序列 + `make check-stub` 数字验收 |
| V 简单优先 | ✅ | 不新建 service 域（D1 替代方案已评估）；楼层 config 不做 schema 校验（YAGNI）；helper 复用 |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD；B2C 定位（装修是平台内容, 无商户维度） |
| 跨域协作 | ✅ | 零跨域（商品装配在 shop 域内） |

**Phase 1 复查**：✅ 四产物无越界；D1（新建接口）与 D2（config 约定）已在 research 记录替代方案与理由。

## Project Structure

### Documentation (this feature)

```text
specs/009-logistics-ops/
├── plan.md / research.md / data-model.md / quickstart.md
├── contracts/operation-permissions.md
├── checklists/requirements.md
└── tasks.md（/speckit-tasks 后续生成）
```

### Source Code (repository root)

```text
apps/server/
├── internal/
│   ├── service/shop/
│   │   ├── misc.go                    # ILogisticsLogic 微扩 Detail; 新建 IOperationLogic（D1/D7）
│   │   ├── logistics_impl.go          # 物流公司 5 方法（新增）
│   │   ├── operation_impl.go          # 装修 10 方法（新增）
│   │   └── operation_impl_test.go     # TDD 红绿（新增）
│   ├── model/dto_shop.go              # 新增 7 个 DTO（D1）
│   ├── errcode/errcode.go             # 新增 40012 物流编码已存在（D5）
│   ├── controller/admin/              # logistics ×5 + banner ×4 + floor ×4 桩填充
│   └── controller/shop/               # banner_list + floor_list 桩填充
└── specs/PROGRESS.md                  # 批次状态（同 commit 更新）
```

**Structure Decision**: 全落既有分层与既有包；新增 2 个实现文件（logistics_impl/operation_impl）+ 测试。

## Complexity Tracking

> 无违例——装修域接口/DTO 新建属批次计划内工作（PROGRESS 批次 03 行已标注"装修接口/DTO 需新建"）。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
