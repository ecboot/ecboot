# Quickstart: 011-member-center 验证指南

## 前置

```bash
cd apps/server
docker compose up -d
make migrate-up               # 至 000036（含本批偏好表）
make migrate-fresh            # 空库全量重放零失败（SC-004）
```

## 自动化验证（DoD 主通道）

```bash
make test                     # 全绿（含既有认证/交易/商品回归 + 本批会员测试）
make lint                     # 本批文件零问题
make check-stub               # user 桩 34→13（本批 21 清零）
```

## 手工验证序列（可选，curl；会员令牌取自短信登录）

1. **资料**：`GET /user/profile` → 手机号**脱敏**（如 138****1111）、等级名与成长值；
   `PUT /user/profile` {nickname} → 再查已变更。
2. **登录记录**：`GET /user/login-logs` → 含刚登录的一条。
3. **地址**：建两条 → `PUT /user/addresses/{id}/default` 设 B → 查列表仅 B 为默认；
   用**他人地址 ID** 改/删 → 10006（越权防护）。
4. **收藏**：`PUT /user/favorites/{spuId}` → 列表含实时价态 → DELETE（软删）→ 再 PUT（**复活**，无重复行）。
5. **足迹**：浏览后 `GET /user/footprints` → `DELETE /user/footprints` 清空。
6. **消息**：造 2 未读 → `GET /user/messages` 未读计数=2 → 单条已读 → 全读 → 计数=0。
7. **偏好**：`GET /user/notify-preferences` 默认全开 → `PUT` 关闭短信 → 再查为关闭。
8. **积分**：`GET /user/point-account` 与 `/user/point-logs?bizType=` 正常（余额可为负）。
9. **邀请**：`GET /user/distribution/invite-records` 仅本人记录。

## 验收红线（宪法 IV）

任何"完成"声明前，`make test` 与 lint 的实际输出必须为通过；桩数以 `make check-stub` 输出为准。
