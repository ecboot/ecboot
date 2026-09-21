# Quickstart: 商品目录域验证指南

前提：`cd apps/server && docker compose up -d && make migrate-up`（含 000034）+ 后台种子超管可用（管理端接口需登录态——过渡期形态见 research D4）。

## 场景一：后台建数据 → C 端可见（SC-002 正向）

```bash
BASE=http://127.0.0.1:8080
# 管理端（以种子超管登录, 具体 follow 后台登录实现形态）
# ① 分类/品牌
# ② 创建 SPU + 2 个启用 SKU + 库存调整 +100
# ③ SPU 上架
curl -s "$BASE/shop/categories" | head -c 120   # 分类树含新建节点
curl -s "$BASE/shop/products?categoryId=<id>"   # 列表含该商品, 价格区间=SKU min~max
curl -s "$BASE/shop/products/<spuId>"           # 详情: 2 SKU 可售 + 评价汇总(avg 为 "—")
```

## 场景二：管理规则正反（SC-003）

| 用例 | 操作 | 预期 |
|---|---|---|
| 规格冲突 | 同 SPU 再加相同规格组合 SKU | 30006 |
| 无启用 SKU 上架 | 全部 SKU 禁用后上架 | 30005 |
| 分类引用禁删 | 删除有商品的分类 | 30004 |
| 库存负向超调 | 可售 10 调整 -20 | 30007 |
| 调整留痕 | 调整 +50 | inventory_log +1 条, 可售 60 |

## 场景三：下架/删除反向（SC-002 反向）

SPU 下架 → C 端列表/详情不再返回（详情 30001 语义）；删除分类后其商品在列表仍可见但分类树无该节点。

## 场景四：足迹（US1 FR-006）

会员凭证查看详情两次 → user_footprint 该 (user, spu) 仅一行、view_count/last_view_at 更新。

## 宪法 IV

`cd apps/server && make build && make test`（service 层商品/库存测试全绿为最终标准）。
