# Implementation Plan: 认证引导纵切片（auth-bootstrap，Go 版）

**Branch**: `003-auth-bootstrap` | **Date**: 2026-09-20 | **Spec**: [spec.md](./spec.md)

## Summary

打通首个业务纵切片（Go 版）：统一响应/全局异常基座（common+web）→ 图形/短信验证码组件（infra，短信 mock）→ 验证码端点（api/commonc）→ 手机号注册即登录/微信归并/休眠核身（service/user + api/user）→ 会话凭证（opaque token + gredis）。**新迁移仅 1 个**（000032：7 条 system_config 种子），70 表结构零变更。

## Technical Context

**Language/Version**: Go 1.25（以 go.mod 为准）+ GoFrame v2.10.3

**Primary Dependencies**: 栈内（goframe：web/config/gdb/gredis/gtest）；新增 1 个半官方库 `golang.org/x/image`（图形验证码字体，自绘实现）

**Storage**: MySQL 8.4（gdb；dao/model 由 `make gen` 生成）+ Redis（gredis：验证码/会话/频控，键模型见 data-model.md）；迁移 golang-migrate（`apps/server/migrations`）

**Testing**: GoFrame `gtest` + 标准 `testing`；端到端用 `gtest` 起 HTTP 进程内探测（或 httptest）打全链路（验证码→登录→me→logout）；宪法 IV 命令 `make test`（go test ./...）

**Target Platform**: apps/server 各 internal 包（落位见结构树）；运行依赖 compose 全栈

**Project Type**: Go 模块化单体纵切片（REST API 后端）

**Performance Goals**: 登录链路三接口 p95 < 500ms（本地容器口径）

**Constraints**: 分层方向遵循工程既有 GoFrame 结构（api 定义 → controller → service → repository/dao）；手机号库内零明文；阈值全走 `system_config`（30 分钟生效）；B2C 定位（无商户/商家概念）；生成物勿手改

**Scale/Scope**: 新增约 25 个 Go 文件、8 端点、1 个迁移；六个 internal 包首次落业务代码

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.0.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位既有 GoFrame 分层包；分层方向遵循工程约定（api 定义 → controller → service → repository/dao） |
| II 统一技术栈 | ✅ | Go+GoFrame 栈内组件（gdb/gredis/gtest）；唯一新依赖 `golang.org/x/image`（半官方，验证码字体）已说明用途 |
| III 中文优先 | ✅ | 文档/注释中文；标识符英文 |
| IV 可验证交付 | ✅ | `make test`/`make build`/`make lint` 为证；dao/model 生成物勿手改；quickstart 可复制命令 |
| V 简单优先 | ✅ | 70 表零变更（仅 1 配置种子迁移）；opaque token 而非 JWT；会话/验证码收敛 gredis |
| 工程约束 | ✅ | 迁移位置/格式合规；B2C + 多门店定位：用户域无门店/商户交互面（显式声明） |

**Phase 1 复查**：✅ research/data-model/contracts/quickstart 无越界。

## Project Structure

### Documentation (this feature)

```text
specs/003-auth-bootstrap/{plan,research,data-model,quickstart}.md, contracts/api-contracts.md, tasks.md(后续)
```

### Source Code (apps/server)

```text
api/
├── common/v1/captcha.go      验证码接口定义（image/verify/sms + mock 取码，条件注册）
└── user/v1/auth.go           登录接口定义（sms-login/wx-login/logout/me）
internal/
├── controller/{common,user}  控制器（参数绑定→service 调用→响应组装）
├── service/user              业务接口 + 实现（注册即登录/微信归并/休眠核身决策树；WxClient mock）
├── repository/               数据访问封装（gf gen dao 之上）
├── middleware/               认证（Bearer token 校验/续期）与统一响应/恢复
├── library/                  技术组件（验证码/短信 mock/手机号加密/会话）实现时定
└── {app,bootstrap,routes,...} 既有装配
migrations/000032_auth_config.up.sql   7 条 system_config 种子
```

**Structure Decision**: 遵循工程既有 GoFrame 分层（api 定义 / controller / service / repository）；业务逻辑集中 service/user，控制器薄组装；跨包依赖经接口（WxClient 在领域内定义，mock 同包）。

## Complexity Tracking

> 无宪法违例。边界项：mock 取码端点条件注册（仅开发态配置开启），已在 contracts 标注。
