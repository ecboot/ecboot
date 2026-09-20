# ECBOOT 后端（apps/server）

社交电商平台后端——**Go + GoFrame v2** 的模块化单体（宪法 2.0.0），由 `main.go` 装配为单一部署单元。

## 模块地图

```
main.go ── 装配四渠道（HTTP API 的唯一宿主）
├─ internal/api/user    小程序前台 · 会员中心 API
├─ internal/api/shop    小程序前台 · 商城 API
├─ internal/api/admin   管理后台 API
├─ internal/api/commonc 公共 API（短信/图形验证码等）
│     └── 四渠道共同依赖 → internal/web（Web 基座：统一响应/异常/认证/TraceId）
│                          └─→ internal/common / internal/infra（基础层）
└─ internal/service/user / internal/service/shop（领域服务，DDD 内部分层）
internal/dao + internal/model = gf gen dao 生成物（勿手改，宪法 IV）
```

## 分层与强制规则

- 依赖方向：`main → api 渠道 → service 领域 → common/infra`，**四渠道互不 import**、领域兄弟互禁——由 `.golangci.yml` 的 **depguard** 机器强制（Go 版 enforcer），违规即 lint 失败。矩阵权威：`specs/001-maven-module-deps/contracts/module-contracts.md`。
- 各包有 doc.go 说明定位；领域包内部 DDD 分层随特性 003 起展开。

## 构建与运行

```bash
cd apps/server
docker compose up -d                # MySQL 8.4(宿主13306) / Redis / Elasticsearch
make build                          # = go build ./...（宪法 IV）
make test                           # = go test ./...
make lint                           # golangci-lint（含 depguard 分层强制）
make run                            # 启动（:8080，配置 manifest/config/config.yaml）
make gen                            # gf gen dao（连接库生成 dao/model）
make migrate-up                     # golang-migrate 应用迁移
make migrate-fresh                  # 空库全量重放验证
```

## 数据库

- 迁移：`migrations/`（golang-migrate 格式 `NNNNNN_name.up/down.sql`，000001~000031，**70 张业务表**；历史 Flyway V1~V31 已一次性转换，语义零变更）
- Schema 设计文档：`docs/schema-design.md`
- ORM：GoFrame `gf gen dao` 生成物（`internal/dao`、`internal/model`），勿手改

## 关键约束（宪法 2.0.0 摘要）

- 中文 Conventional Commits；文档中文、标识符英文
- 金额一律 `DECIMAL(10,2)` + decimal 库等价精度处理，禁止 float/==
- 合规红线：分销≤2 级（结构强制）、注销匿名化、敏感字段加密+哈希索引
- **多商户演进预留（B2C 严守）**：设计留多商户之门（seller 维度），功能守单商户
- 完整宪法：`.specify/memory/constitution.md`
