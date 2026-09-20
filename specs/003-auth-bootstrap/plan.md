# Implementation Plan: 认证引导纵切片（auth-bootstrap，Go 版）

**Branch**: `003-auth-bootstrap` | **Date**: 2026-09-20 | **Spec**: [spec.md](./spec.md)

> 技术栈变更（宪法 2.0.0）：本 plan 为 Go + GoFrame v2 版本，替代 2026-09-18 的 Java 版技术方案（业务 spec 不变）。

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

**Constraints**: depguard 分层矩阵（.golangci.yml）；手机号库内零明文；阈值全走 `system_config`（30 分钟生效）；B2C 严守（无 seller 交互面）；生成物勿手改

**Scale/Scope**: 新增约 25 个 Go 文件、8 端点、1 个迁移；六个 internal 包首次落业务代码

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则/约束（2.0.0） | 结论 | 依据 |
|---|---|---|
| I 模块化单体 | ✅ | 全部落位既有 internal 包；分层方向/渠道互禁/兄弟互禁由 depguard 机器强制（矩阵已配置） |
| II 统一技术栈 | ✅ | Go+GoFrame 栈内组件（gdb/gredis/gtest）；唯一新依赖 `golang.org/x/image`（半官方，验证码字体）已说明用途 |
| III 中文优先 | ✅ | 文档/注释中文；标识符英文 |
| IV 可验证交付 | ✅ | `make test`/`make build`/`make lint` 为证；dao/model 生成物勿手改；quickstart 可复制命令 |
| V 简单优先 | ✅ | 70 表零变更（仅 1 配置种子迁移）；opaque token 而非 JWT；会话/验证码收敛 gredis |
| 工程约束 | ✅ | 迁移位置/格式合规；多商户预留（B2C 严守）：用户域无 seller 交互面（显式声明） |

**Phase 1 复查**：✅ research/data-model/contracts/quickstart 无越界。

## Project Structure

### Documentation (this feature)

```text
specs/003-auth-bootstrap/{plan,research,data-model,quickstart}.md, contracts/api-contracts.md, tasks.md(后续)
```

### Source Code (apps/server)

```text
internal/
├── common/           api.go(Response/错误码) exception.go paging.go
├── infra/
│   ├── captcha/      service.go(生成/校验/频控) image.go(自绘) ticket.go
│   ├── sms/          sender.go(接口) mock_sender.go
│   ├── security/     phone_cipher.go(AES-GCM+盐哈希) session.go(token/续期/登出) keys.go
│   └── config/       sysconfig.go(system_config 读取,缓存≤30min)
├── web/
│   ├── response.go   统一响应写出
│   ├── exception.go  全局恢复/错误映射中间件
│   ├── traceid.go    追踪号中间件
│   └── auth.go       Bearer 认证中间件 + CurrentUser 上下文注入
├── api/
│   ├── commonc/captcha_handler.go   (+mock 取码路由,条件注册)
│   └── user/auth_handler.go
├── service/user/
│   ├── auth.go       应用服务(注册即登录/归并/休眠核身决策树)
│   ├── wx.go         WxClient 接口 + mock 实现(配置开关)
│   └── repo.go       仓储(gf gen dao 之上的领域封装)
main.go                路由装配 + 中间件挂载
migrations/000032_auth_config.up.sql   7 条 system_config 种子
```

**Structure Decision**: 领域逻辑集中 `service/user`（DDD 细分随域成长展开）；渠道 handler 只做参数绑定与组装；跨包依赖经接口（WxClient 在领域内定义，mock 同包）。

## Complexity Tracking

> 无宪法违例。边界项：mock 取码端点条件注册（仅开发态配置开启），已在 contracts 标注。
