# Quickstart: 008-store 验证指南

> 目标：不读实现代码即可验证批次 DoD。实现细节见 [plan.md](./plan.md) 与 tasks.md。

## 前置

```bash
cd apps/server
docker compose up -d          # MySQL 13306 / Redis
make migrate-up               # 至 000035（本批零新迁移）
```

## 自动化验证（DoD 主通道）

```bash
make test                     # 全绿（含批次 01 认证/权限回归 + 本批 store 测试）
make lint                     # 本批文件零问题
make check-stub               # admin 桩 109→104、common 桩 3→1（本批 7 清零）
```

## 手工验证序列（可选，curl）

1. **创建门店**（超管）：`POST /admin/stores`（名称/A 省市区划码/地址/坐标/自提开）→ 返回系统生成的
   `storeNo`（ST 前缀）；再建一家 → 编码不同。
2. **游客区县检索**：`GET /common/stores?districtCode=<A区码>` → 仅 A 区营业店。
3. **游客附近检索**：`GET /common/stores?longitude=&latitude=` → 半径内门店按距离升序、带 distanceM；
   歇业店与无坐标店不出现；`radiusKm` 扩大后命中变多。
4. **游客详情**：`GET /common/stores/{id}` → 完整档案 + 状态 + 自提标识；不存在 ID → 10006。
5. **后台管理**：`PUT /admin/stores/{id}` 置歇业 → 游客列表消失、详情可见（status=2）；
   `GET /admin/stores?keyword=<编码>` 命中；`DELETE /admin/stores/{id}` 后全端点不可见。
6. **权限拦截**：无权账号调 `POST /admin/stores` → 10005；超管成功（延续批次 01 机制）。

## 验收红线（宪法 IV）

任何"完成"声明前，`make test` 与 lint 的实际输出必须为通过；桩数以 `make check-stub` 输出为准。
