# Quickstart: 009-logistics-ops 验证指南

> 目标：不读实现代码即可验证批次 DoD。

## 前置

```bash
cd apps/server
docker compose up -d
make migrate-up               # 至 000035（本批零新迁移）
```

## 自动化验证（DoD 主通道）

```bash
make test                     # 全绿（含批次 01/02 回归 + 本批 logistics/operation 测试）
make lint                     # 本批文件零问题
make check-stub               # admin 桩 104→91、shop 桩 42→40（本批 15 清零）
```

## 手工验证序列（可选，curl；中文参数请用 URL 编码或 --data-binary）

1. **建物流公司**（超管）：`POST /admin/logistics-companies` {code:"SF", name:"SF Express"} → 返回 ID；
   同 code 再建 → 40012；`GET …/{id}` 详情可见。
2. **停用保留**：`PUT …/{id}` {status:0} → 列表（status=0 筛选）仍可检出、详情可见。
3. **软删**：`DELETE …/{id}` → 列表与详情均 10006。
4. **建轮播**：`POST /admin/banners` {position:1, imageUrl:"…", startTime:"", endTime:""} → 返回 ID；
   `GET /admin/banners?position=1` 命中。
5. **投放语义**：再建两条（startTime 设为未来 / endTime 设为过去）→ `GET /shop/banners?position=1` 
   **仅返回长期那条**（未到时段与已过期被排除）。
6. **建楼层**：`POST /admin/floors` {floorType:2, title:"热销", config:{"spuIds":["<有效spuId>","999999"]}} →
   `GET /shop/floors` 返回该楼层且 Products 仅含有效商品（失效 ID 被剔除不报错）。
7. **权限**：无权账号 `POST /admin/banners` → 10005；超管成功；持 banner 权限不持 floor 权限者建楼层被拒。
8. **公开性**：无 token 调 `GET /shop/banners?position=1` 与 `GET /shop/floors` → 正常返回。

## 验收红线（宪法 IV）

任何"完成"声明前，`make test` 与 lint 的实际输出必须为通过；桩数以 `make check-stub` 输出为准。
