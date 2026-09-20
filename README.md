# ECBOOT

B2C + 多门店的社交电商系统（小程序商城 + 管理后台）。

- 后端：Go + GoFrame v2（`apps/server`，模块化单体）
- 前端：TanStack Start（`apps/web` 商城 / `apps/admin` 后台）、uni-app（`apps/mobile`）
- 数据库：MySQL 8.4（70 表，golang-migrate 迁移 `apps/server/migrations`）
- 设计文档：`docs/schema-design.md`（Schema）、`docs/layer-contracts.md`（分层）、`docs/adr/`（决策）
