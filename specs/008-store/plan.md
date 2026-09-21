# Implementation Plan: 门店域（008-store）

**Branch**: `008-store` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

实现线下门店域（自提/核销载体）：7 端点（admin 管理 5 + common 游客检索 2）的 service 实现与 controller 连线。
关键技术点：附近检索 = 包围盒预筛 + Haversine 计算列（dao 链内, D1）；门店编码 sonyflake 首次接线
（新叶子包 library/idgen, D2）；AdminDetail 接口微扩（D3）。零表结构变更、零新 DTO、零新错误码。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3（gdb + dao 生成物 + do 强类型）

**Primary Dependencies**: 既有栈复用；`github.com/sony/sonyflake/v2`（已在 go.mod, 本批提升为直接使用）

**Storage**: MySQL 8.4（store 表既有, 000031；idx_location/idx_district_status 命中）；零新迁移

**Testing**: `internal/service/shop/store_impl_test.go` gtest 红绿（确定性配置注入 + 数据自建清理；
该包已有测试 init 基座，直接复用）

**Target Platform**: `internal/service/shop`（store_impl）、`internal/library/idgen`（新叶子包）、
controller `admin_v1_admin_store_*.go` ×5 + `common_v1_store_*.go` ×2

**Project Type**: GoFrame 模块化单体·门店域纵切片

**Performance Goals**: 门店量级数百，无专项指标；附近检索走 idx_location 包围盒

**Constraints**: 分层契约 §二/§四（do/entity 强类型、controller 禁触 dao）；TDD 红绿；
宪法 V（idgen 为 library 既有层内新叶子包——与 captcha/sms 同级, 非新顶层包）

**Scale/Scope**: 7 端点、service 包级函数 7 个（含 AdminDetail 微扩）、距离球面公式 1 处、测试覆盖全部验收场景

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | service/shop 既有域 + 既有 controller 渠道；跨渠道零引用 |
| II 统一技术栈 | ✅ | sonyflake 为约定内既有依赖（go.mod 已在, 提升直接使用, AGENTS.md 基础库约定） |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | TDD 红绿 + quickstart 序列 + `make check-stub` 数字验收 |
| V 简单优先 | ✅ | 不建顶层包（idgen 系 library 层内叶子包, 同 captcha/sms 形态）；repository 未接线不启用；零迁移零 DTO 零错误码 |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD；B2C 多门店定位（store 即线下载体, 无商户概念） |
| 跨域协作 | ✅ | 零跨域（门店域自洽；交易侧消费归 012+） |

**Phase 1 复查**：✅ 四产物无越界；D1/D2 两个技术决策已在 research 记录替代方案与理由。

## Project Structure

### Documentation (this feature)

```text
specs/008-store/
├── plan.md                        # 本文件
├── research.md                    # Phase 0：D1~D5 决策与勘察
├── data-model.md                  # Phase 1：实体/状态机/附近检索数据流
├── contracts/store-permissions.md # Phase 1：7 端点 × 权限点映射
├── quickstart.md                  # Phase 1：验证指南
├── checklists/requirements.md     # specify 阶段质量单
└── tasks.md                       # Phase 2（/speckit-tasks 生成，本批后续）
```

### Source Code (repository root)

```text
apps/server/
├── internal/
│   ├── library/idgen/idgen.go       # sonyflake 单例 + NextID（新增叶子包, D2）
│   ├── service/shop/
│   │   ├── misc.go                  # IStoreLogic 接口微扩 AdminDetail（D3）
│   │   ├── store_impl.go            # 7 方法实现（新增）
│   │   └── store_impl_test.go       # TDD 红绿（新增）
│   ├── controller/admin/            # store_list/create/update/delete/detail 桩填充 ×5
│   └── controller/common/           # store_list/store_detail 桩填充 ×2
└── specs/PROGRESS.md                # 批次状态（同 commit 更新）
```

**Structure Decision**: 全部落位既有分层与既有包（宪法 V）；唯一新文件为
`library/idgen/idgen.go`（约定首次接线）、`store_impl.go` + 测试。

## Complexity Tracking

> 无违例——idgen 叶子包系 library 既有层内的自然扩展（同 captcha/sms 形态），不构成新顶层包。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
