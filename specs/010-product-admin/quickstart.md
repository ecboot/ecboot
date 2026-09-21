# Quickstart: 010-product-admin 验证指南

> 目标：不读实现代码即可验证批次 DoD。

## 前置

```bash
cd apps/server
docker compose up -d
make migrate-up               # 至 000035（本批零新迁移）
```

## 自动化验证（DoD 主通道）

```bash
make test                     # 全绿（含 005 商品测试、006 交易链路回归 + 本批库存测试）
make lint                     # 本批文件零问题
make check-stub               # admin 桩 91→69、shop 桩 40→35（本批 27 清零）
```

## 手工验证序列（可选，curl）

1. **后台目录**：`GET /admin/categories` 树 → `POST /admin/categories` 建子类 → 树含新节点；
   `POST /admin/brands` 建品牌 → `GET /admin/brands` 可见。
2. **SPU/SKU**：`POST /admin/products`（含 SKU 与库存）→ 返回 id 与编码 → 详情可见 →
   `POST /admin/products/{id}/status` 状态切换生效。
3. **库存**：`GET /admin/inventories` 可见售推导 → `POST /admin/inventories/{skuId}/adjust` {delta:5}
   → 总量 +5 且 `inventory_log` 新增 change_type=5 一行（operator=admin:{id}）；
   `GET /admin/inventories/warnings` 按阈值命中。
4. **库存不可为负**：调整 −999 → 30007。
5. **C 端浏览**（无凭证）：`/shop/categories`、`/shop/brands`、`/shop/products`、`/shop/products/{id}`、
   `/shop/search?keyword=` 正常；下架商品不出现在列表与搜索。
6. **权限**：无权账号调 `POST /admin/products` → 10005；超管成功。

## 验收红线（宪法 IV）

任何"完成"声明前，`make test` 与 lint 的实际输出必须为通过；桩数以 `make check-stub` 输出为准。
