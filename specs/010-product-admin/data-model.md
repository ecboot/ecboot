# Data Model: 010-product-admin

> 零表结构新增；库存三方法为本批新写，其余端点复用既有实现。

## 一、库存语义（本批新写部分）

### inventory（SKU 维度）
| 字段 | 语义 | 本批规则 |
|---|---|---|
| sku_id | 主键（与 SKU 1:1） | 查询/调整的键 |
| total | 总库存 = 可售 + 锁定 | 调整目标列（有符号增减） |
| locked | 下单锁定（未支付） | 只读展示；交易侧维护 |
| warn_count | 预警阈值 | 预警判定依据 |

**约束**：表级 `CHECK (locked <= total)`；调整 MUST NOT 使 total 为负（应用层条件更新保证）。

### inventory_log（只追加流水）
| 字段 | 本批写入 |
|---|---|
| change_type | **5（后台调整）** |
| quantity | `|delta|`（绝对值, 方向由 change_type 定） |
| total_after / locked_after | 调整后快照（对账用） |
| operator | **`admin:{id}`**（表注释约定格式） |
| remark | 运营备注 |

## 二、库存三方法数据流

```
InventoryList(q)     ──▶ inventory JOIN product_sku（sku_no/sku_name）
                          available = total - locked（推导, 不落库）
InventoryWarnings(p) ──▶ WHERE total - locked <= warn_count（SQL 侧）
InventoryAdjust(skuId, delta, remark)
  └─ 事务
       ├─ UPDATE inventory SET total=total+? WHERE sku_id=? AND total+?>=0
       │    affected=0 → 30007（不存在或将为负）
       └─ INSERT inventory_log(change_type=5, quantity=|delta|, total_after, locked_after,
                               operator=admin:{id}, remark)
```

## 三、既有实体（本批仅连线, 不改语义）

- **分类/品牌/SPU/SKU**：字段与规则见 005 规格；端点映射见 contracts。
- **C 端浏览口径**：价格区间文本（price_min/price_max）、可售标记、首图——由既有实现决定。

## 四、端点 × 方法 × 权限

见 [contracts/product-endpoint-mapping.md](./contracts/product-endpoint-mapping.md)（27 端点全表）。

## 五、DTO 与错误码

- **DTO 零新增**（既有 76 类型覆盖）。
- **错误码零新增**（复用 30001~30007 与 10001/10005/10006）。
