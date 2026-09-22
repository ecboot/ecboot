# Feature Specification: 会员管理与审计风控（批次12）

**Feature Branch**: `018-member-audit-risk` | **Created**: 2026-09-22 | **Status**: Draft

**Input**: 批次 12：**13 个 admin 端点**（全为桩）——**会员管理 4**（列表/详情/禁用启用/改绑手机号）、**风控 7**（规则 CRUD+详情、事件列表、申诉处理）、**审计日志 2**（登录日志/操作日志列表）。接口现状：`IAuditLogic`（审计侧）已有待实现；会员管理与风控管理接口**缺失需新建**（D3 记账）。权限点 7 个全在 000032 种子（member:read/update、risk:rule:read/manage、risk:record:read/appeal、system:audit:read），**预期零权限迁移**。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 会员管理（Priority: P1）🎯 MVP

运营者按关键词/状态筛选会员分页列表，查看会员详情（资料+等级+统计概要），**禁用/启用**账号（禁用即登录与业务操作被拒），以及**改绑手机号**（实名场景的客服通道；新手机号唯一性校验）。

**Why this priority**: 会员管控是后台对 C 端账号的唯一治理手段；禁用语义横跨登录与全部业务域。

**Independent Test**: 造会员 → 列表可见 → 禁用 → 会员侧登录/操作被拒（status 语义既有）→ 启用恢复 → 改绑新手机号 → 详情反映新号 → 旧号可重新注册。

**Acceptance Scenarios**:

1. **Given** 既有会员，**When** 列表按关键词（昵称/手机号）查询，**Then** 分页返回且手机号脱敏展示。
2. **Given** 正常会员，**When** 禁用，**Then** `user.status` 置禁用值（条件更新判行数）；该会员后续登录被既有认证口径拒绝。
3. **Given** 已禁用会员，**When** 启用，**Then** 恢复正常。
4. **Given** 改绑手机号, **When** 新号已被占用，**Then** 拒绝（`uk_phone_hash` 唯一键 → 业务码）；**When** 新号合法，**Then** phone/phone_hash 原子更新且详情回读一致。
5. **Given** 不存在的 userId，**Then** 详情/禁用/改绑均返回业务码 404 类，不泄露存在性差异之外的内部信息。

---

### User Story 2 - 风控规则与事件（Priority: P1）

运营者管理风控规则（规则类型 1黑名单 2高频下单 3异常领券 4佣金套利 5休眠分级；动作 1拦截 2标记；条件表达式与启停），查看风控事件列表（命中记录），处理**申诉**（通过=解除/驳回=维持，写申诉结论）。规则引擎评估器（`IRiskHit` 端口实现）按启用规则对入口调用给出拦截/放行判定——批次 09/11 预留的端口消费点就此闭合。

**Why this priority**: 风控是被刷重灾区的最后防线（000029 schema 级要求）；管理面与评估器同批闭合才能让"挂风控"从降级放行变成真实防线。

**Independent Test**: 建"黑名单"规则指向某用户 → 该用户命中玩法入口被拦截 → 风控事件落记录（action=拦截）→ 申诉通过 → 记录申诉状态更新且后续放行。

**Acceptance Scenarios**:

1. **Given** 合法规则参数（type/action/condition_expr），**When** 创建，**Then** 规则可见且详情回读一致；作用域非法 → 拒绝。
2. **Given** 启用的黑名单规则（condition 指向用户 U），**When** U 调用风控入口（`RiskHit.Hit`），**Then** 返回拦截=true 且落 risk_record（action=1）。
3. **Given** 停用/软删规则，**When** 评估，**Then** 不生效。
4. **Given** 风控事件，**When** 申诉通过，**Then** `appeal_status` 置通过（条件更新判行数）且事件关联的用户不再被该记录拦截。
5. **Given** 重复申诉（已有结论），**When** 再次处理，**Then** 拒绝。
6. **Given** 规则软删, **Then** 列表消失；规则不存在 → 404 类业务码。

---

### User Story 3 - 审计日志查询（Priority: P2）

后台查看**登录日志**（管理员登录成功/失败、IP、UA）与**操作日志**（写操作审计——批次 01 AOP 写入侧已产出数据），支持筛选分页。本批只做**查询面**（写入侧已在线）。

**Why this priority**: 纯查询；数据已由批次 01 的审计写入持续产出，补齐可观测闭环。

**Independent Test**: 造登录/操作日志行 → 列表分页可见、筛选正确。

**Acceptance Scenarios**:

1. **Given** 既有登录日志，**When** 分页查询（按管理员/结果筛选），**Then** 返回正确且敏感字段按既有脱敏口径。
2. **Given** 既有操作日志（批次 01 AOP 产出），**When** 按模块/时间筛选，**Then** 分页正确。

---

### Edge Cases

- **禁用语义一致性**: 禁用只改 `user.status`（既有认证/业务域已消费该列），不级联清理会话/数据（与批次 01 登录口径对齐）。
- **改绑手机号并发**: 同新号并发改绑 → 唯一键兜底 1062 → 业务码。
- **申诉幂等**: 已有结论的记录重复申诉 → 拒绝（条件更新）。
- **规则评估无规则/无记录**: 空规则集 → 放行（不阻断业务）；评估失败不阻断主流程（记录错误告警）。
- **分页越界**: 空列表 + 正确 total。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-1（会员）**: 列表（关键词/状态筛选+分页+手机号脱敏）、详情、禁用/启用（条件更新）、改绑手机号（唯一键兜底 1062→业务码）。
- **FR-2（规则）**: 规则 CRUD+详情；类型/动作/条件校验；软删；命中评估器按启用规则判定（黑名单类型 V1 全量，表达式引擎 V1 按 condition_expr 的结构化字段判定——口径实现期固化注释）。
- **FR-3（事件与申诉）**: 命中落 risk_record（用户/规则/动作/上下文）；申诉处理条件更新（appeal_status 0无 1申诉中 2通过 3驳回，已有结论拒绝）。
- **FR-4（审计）**: 登录日志/操作日志查询面（IAuditLogic 已有方法实现），筛选分页。
- **FR-5（权限）**: member:read/update、risk:rule:read/manage、risk:record:read/appeal、system:audit:read——按路由注释挂接（RequirePerm 首行），权限点已就位。
- **FR-6（风控接线）**: `IRiskHit` 端口在本批由规则评估器实现并 bootstrap 装配（批次 09/11 的降级放行点转为真实判定）；评估失败不阻断主流程。

## Key Entities

- `user`（status 禁用语义/phone_hash 唯一键）、`risk_rule`（000021：type/action/condition_expr/status/deleted）、`risk_record`（000021：user/rule/action/appeal_status）、`admin_login_log`/`admin_operation_log`（批次 01 既有）。

## Success Criteria *(mandatory)*

- **SC-1**: `check-stub` 对账 **admin 16→3**（四渠道合计 16→3；余 3 为批次 13 dashboard）。
- **SC-2**: 全量 `go test ./...` 两连跑全绿；`golangci-lint run` **0 issues**。
- **SC-3**: 风控闭环测试：规则→命中拦截→落记录→申诉通过→放行（全链一条龙）。
- **SC-4**: 零迁移、零权限种子（勘察已确认）；权限挂载经非超管 10005 对照测试（批次 10 I4 内化）。

## Assumptions

- 规则评估器 V1 支持**黑名单类型**（condition_expr 指向用户/维度）+ 通用类型按记录落库；复杂表达式引擎（高频/套利的实时统计判定）按 spec 边界记为"记录与申诉面完整、实时统计判定随部署面数据积累后续增强"——**记账**。
- 风控入口消费点（玩法/领券）已在批次 09/11 预留 `RiskHit` 调用；本批装配后即生效。
- 申诉通过后"解除拦截"的落地口径：该记录状态更新（历史记录不再驱动新判定）；黑名单规则的持续拦截由规则本身存在与否决定（记录不直接驱动拦截）。
