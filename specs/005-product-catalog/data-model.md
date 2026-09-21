# Data Model: 商品目录域（Phase 1）

**零新表**；1 个迁移（000034 SPU 价格冗余列）。数据形态 = 既有表映射（gf dao）+ 冗余列刷新规则。

## 一、迁移 000034_spu_price_range

```sql
ALTER TABLE `product_spu`
  ADD COLUMN `price_min` DECIMAL(10,2) NULL COMMENT '启用SKU最低价(冗余,SKU变更时刷新;全部禁用为NULL)' AFTER `sale_count`,
  ADD COLUMN `price_max` DECIMAL(10,2) NULL COMMENT '启用SKU最高价(冗余)' AFTER `price_min`,
  ADD INDEX `idx_price_min` (`price_min`);
```

- 刷新触发点：SKU 创建/编辑/启停/删除（同事务重算该 SPU 全部启用 SKU 的 min/max）。
- 列表价格区间筛选/价格排序直接基于本列（闭区间语义在应用层拼接条件）。

## 二、读模型（entity → model DTO 转换口径）

| DTO | 来源与口径 |
|---|---|
| ProductCard | spu（上架+未删）+ price 冗余列区间；`SaleCount` 排序字段 |
| ProductDetail.Skus | 仅启用 SKU；`Sellable` = SPU 上架 AND SKU 启用 AND 可售>0（联查 inventory） |
| ReviewSummary | product_review 按 `audit_status=1` 聚合（AVG/分档 COUNT）；无评价时 Avg="—" |
| AdminSkuDetail | 全量 SKU（含禁用），含 cost_price（管理可见，C 端 Detail 不转换该字段） |

## 三、状态与规则（实现要点）

- **上架校验**：SPU 下状态置 1 前，必须存在 `status=1 AND deleted=0` 的 SKU（否则 30005）。
- **规格组合唯一**：依赖 V23 生成列 `specs_hash` 的 `uk_spu_specs_hash`；INSERT 捕获 1062 → 30006。
- **分类删除**：存在任何（含软删）商品引用 → 30004。
- **库存调整**：`UPDATE inventory SET total=total+delta WHERE sku_id=? AND total+delta>=0`；affected=0 → 30007；同事务写 inventory_log（change_type=5，operator=后台账号）。
- **限售设置**：全量替换 JSON 列（空数组=不限售）。

## 四、排序口径

| sort | 语义 |
|---|---|
| 0 综合 | sale_count DESC, id ASC |
| 1 销量 | sale_count DESC, id ASC |
| 2 价格升 | price_min ASC, id ASC |
| 3 价格降 | price_min DESC, id ASC |
| 4 上新 | created_at DESC, id ASC |

（同值以 id 升序保证分页稳定。）
