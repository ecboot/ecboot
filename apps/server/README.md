# ECBOOT 后端（apps/server）

社交电商平台后端——**Go + GoFrame v2** 的模块化单体（宪法 2.0.0），由 `main.go` 装配为单一部署单元。

## 模块地图（GoFrame 惯用结构）

```
main.go / internal/app ── 装配与启动
├─ api/{common,user,shop,admin}/v1   四渠道接口定义（互不引用）
├─ internal/controller/*             控制器（绑定→service→响应）
├─ internal/service/{user,shop}      领域业务（兄弟域经事件协作）
├─ internal/repository → dao|model   数据访问（gf gen dao 生成物，勿手改）
├─ internal/middleware               认证/统一响应/恢复
└─ internal/library                  技术组件（验证码/短信/加密/会话）
```

## 分层与强制规则

- 分层契约：`docs/layer-contracts.md`——四渠道目录互不引用、领域兄弟互禁（经事件协作）、生成物勿手改。
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

- 迁移：`migrations/`（golang-migrate 格式 `NNNNNN_name.up/down.sql`，000001~000031，**68 张业务表**）
- Schema 设计文档：`docs/schema-design.md`
- ORM：GoFrame `gf gen dao` 生成物（`internal/dao`、`internal/model`），勿手改

## 关键约束（宪法 2.0.0 摘要）

- 中文 Conventional Commits；文档中文、标识符英文
- 金额一律 `DECIMAL(10,2)` + decimal 库等价精度处理，禁止 float/==
- 合规红线：分销≤2 级（结构强制）、注销匿名化、敏感字段加密+哈希索引
- **系统定位**：纯 B2C + 多门店（store=线下载体：自提/核销/附近门店）；无多租户/多商户概念，B2B2C 另立项目
- 完整宪法：`.specify/memory/constitution.md`
