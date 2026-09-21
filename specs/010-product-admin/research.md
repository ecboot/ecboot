# Research: 010-product-admin 商品目录收口

> Phase 0 产出。基于代码勘察（commit 396f5b5 时点）。

## D1 连线形态跟随既有 struct（不引入包级双轨）

**结论**：controller 经 `shop.NewProductLogic()` 调用既有 struct 方法，与同域既有形态一致；
本批**不**为商品域引入包级函数。

**理由**：shop 域既有实现（`ProductLogicImpl`/`OrderLogicImpl`/`CartLogicImpl`）均为 struct 方法；
批次 01~03 的 system/store/logistics/operation 域用包级函数属另一域。同域内保持单一形态，
避免评审已记录的"两套并存"债务继续扩散。

**替代方案**：新增包级包装函数转发到 struct——凭空多一层且加剧双轨；放弃。

## D2 本批范围 = 连线收口 + 库存补齐（不重写既有）

**结论**：19 个管理方法与 5 个浏览方法保持原样；仅当契约字段与既有签名不匹配时做最小适配
（controller 层转换，或最小改 service 签名并记账）。

**理由**：005 已交付且带测试；重写只引入回归风险。006 的库存锁定/核销以事务内条件更新实现
（`tx.Model("inventory")`），**不动**。

## D3 库存后台新写（三方法, 语义以表注释为准）

**结论**：
- `InventoryList`：按 skuId/关键词筛选分页，返回 {skuId, skuNo, skuName, total, locked, available, warnCount}
  （DTO `model.InventoryItem` 已就位）；`available = total - locked` 推导
- `InventoryWarnings`：SQL 侧 `total - locked <= warn_count` 分页
- `InventoryAdjust`：有符号 delta；条件更新 `WHERE sku_id=? AND total+?>=0`（affected=0 → 30007 既有码）；
  同事务写 `inventory_log`（change_type=5, quantity=|delta|, total_after/locked_after 快照,
  operator=`admin:{id}`, remark）

**理由**：接口 `IInventoryLogic` 已定义前两法与交易侧方法；表注释定义流水字段语义
（operator 格式、change_type 枚举、快照对账用途）。错误码复用既有 30007，零新增。

## D4 预警阈值口径

**结论**：`warn_count` 为 SKU 维度阈值；判定 `available <= warn_count`（含等于, spec FR-007）。

## 勘察结论（非决策）

- **impl 覆盖**：`ProductLogicImpl` 19 管理方法 + 5 浏览方法（均未被 controller 调用）；
  `IInventoryLogic` 无实现者（同 I5 记录的接口债）。
- **DTO 齐备**：`dto_shop.go` 76 个类型含本批所需全部出入参（零新增）。
- **权限码齐备**：product:{category,brand,spu,sku}:{read,create,update,delete} + inventory:{read,adjust}。
- **零迁移**：五表齐备（000002/000003/000034）。
- **测试基座**：009 新增的 `seedFloorSpu`/`ensureRefRow` 可复用于库存测试。
