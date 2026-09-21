# Implementation Plan: 会员中心（011-member-center）

**Branch**: `011-member-center` | **Date**: 2026-09-21 | **Spec**: [spec.md](./spec.md)

## Summary

从零实现 user 域会员自助能力（8 个接口、21 端点）：资料/地址/收藏/足迹/站内信/通知偏好/积分账户/邀请记录。
唯一表结构新增：通知偏好表（`user_notify_preference`）。DTO 基本齐备（`dto_user.go` 23 类型），
仅 invite_record 需补 DTO 与接口方法。

## Technical Context

**Language/Version**: Go 1.25 + GoFrame v2.10.3

**Primary Dependencies**: 既有栈复用（`library/security.PhoneCipher` 脱敏与加密、`library/sms` 等）；零新依赖

**Storage**: MySQL 8.4（user/user_address/user_favorite/user_footprint/user_message/notify_*/point_*/
user_level_rule/user_login_log/invite_record 全既有——**新迁移 000036 建通知偏好表**）

**Testing**: `internal/service/user/*_impl_test.go`（复用该包既有 init 基座与 `issueCode`/`cleanState` 惯例）；
TDD 红绿施于全部新写方法

**Target Platform**: `internal/service/user`（8 个 impl 文件）、`internal/controller/user/*`（21 桩填充）、
`migrations/000036_user_notify_preference.*.sql`

**Project Type**: GoFrame 模块化单体·会员中心纵切片

**Performance Goals**: 无专项（会员自助读多写少）

**Constraints**: 分层契约（controller 只调 service）；**实现形态跟随 user 域既有包级函数**
（`SmsLogin` 先例, 不引入 struct 双轨）；越权防护统一以 userId 条件收口

**Scale/Scope**: 21 端点、8 接口约 20 方法、1 迁移、1 DTO + 1 接口方法微扩、21 controller

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.1.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全落 service/user 既有域与 user 渠道 controller；无跨域引用（同域内用 dao） |
| II 统一技术栈 | ✅ | 零新依赖 |
| III 中文优先 | ✅ | 文档/注释/提交中文 |
| IV 可验证交付 | ✅ | 全部新写方法 TDD；迁移空库重放验证；`make check-stub` 数字 |
| V 简单优先 | ✅ | 不新建域/包；偏好表为最小必要结构（唯一键 user×channel）；未设置=默认全开不预置行 |
| 工程约束 §二/五 | ✅ | do/entity 强类型；TDD |
| 合规（个保法） | ✅ | 手机号脱敏在服务层；明文不出参不出日志（延续既有约束） |

**Phase 1 复查**：✅ 四产物无越界；新增迁移已在 PROGRESS 记账（spec 假设修订）。

## Project Structure

### Documentation (this feature)

```text
specs/011-member-center/
├── plan.md / research.md / data-model.md / quickstart.md
├── contracts/member-endpoint-mapping.md
├── checklists/requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
apps/server/
├── migrations/000036_user_notify_preference.{up,down}.sql   # 新增（会员×渠道唯一）
├── internal/
│   ├── model/
│   │   ├── dto_user.go            # 新增 InviteRecordItem（唯一 DTO 增量）
│   │   └── do/…                   # 生成物（gf gen dao 后）
│   ├── service/user/
│   │   ├── profile_impl.go        # Detail/Update/LoginLogs（+GrowthAdd/LevelRecalc 内部）
│   │   ├── address_impl.go        # 5 端点方法 + GetForOrder（内部）
│   │   ├── favorite_impl.go       # 收藏 3 + 足迹 2（+Record/CleanExpired 内部）
│   │   ├── notify_impl.go         # 消息 3 + 偏好 2（+Enqueue/DispatchTask 内部）
│   │   ├── point_impl.go          # 账户/流水（+Earn/Consume/Refund/ExpireDormant 内部）
│   │   ├── invite_impl.go         # 邀请记录列表（接口微扩）
│   │   └── *_impl_test.go         # TDD 红绿
│   └── controller/user/           # 21 桩填充
└── specs/PROGRESS.md
```

**Structure Decision**: 全部落既有 user 域与渠道；新增 6 个 impl 文件 + 1 迁移 + 1 DTO。

## Complexity Tracking

> 无违例。通知偏好表为 spec FR-015 明确要求的最小结构，已在 PROGRESS 记账。

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| （无） | | |
