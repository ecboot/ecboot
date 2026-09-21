// inventory_impl.go 库存后台（接口契约见 inventory.go IInventoryLogic 的管理面三方法）。
// 语义: 可售 = total - locked（推导, 不落库）; 预警 available <= warn_count（含等于）;
// 调整有符号且不可为负（条件更新, affected=0 → 30007）, 同事务写流水
// （change_type=5 后台调整 / quantity=绝对值 / 快照 / operator=admin:{id}, research D3）。
// 交易侧（锁定/核销/释放/回补）由 006 交易链路实现, 不在本文件范围。
package shop

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// inventoryBaseQuery 库存基础查询（JOIN product_sku 取编码与名称）。
// 注: 不在此设 Fields——Fields 会被 Count 复用生成 `COUNT(col1,col2,...)` 语法错误;
// 取记录时按需指定列。
func inventoryBaseQuery(ctx context.Context) *gdb.Model {
	return dao.Inventory.Ctx(ctx).As("i").
		LeftJoin(dao.ProductSku.Table()+" s", "s.id=i.sku_id")
}

// inventoryCols 取记录时的列清单（含 JOIN 侧编码与名称）。
const inventoryCols = "i.sku_id,i.total,i.locked,i.warn_count,s.sku_no,s.name AS sku_name"

// inventoryFromRecord 行 → DTO（可售推导）。
func inventoryFromRecord(r gdb.Record) model.InventoryItem {
	total := r["total"].Int()
	locked := r["locked"].Int()
	return model.InventoryItem{
		SkuId:     r["sku_id"].Int64(),
		SkuNo:     r["sku_no"].String(),
		SkuName:   r["sku_name"].String(),
		Total:     total,
		Locked:    locked,
		Available: total - locked,
		WarnCount: r["warn_count"].Int(),
	}
}

// InventoryList 库存列表（FR-006）: skuId 精确筛选或关键词（编码/名称）模糊 + 分页。
func InventoryList(ctx context.Context, skuId int64, keyword string, page model.PageReq) (*model.PageResult[model.InventoryItem], error) {
	page = page.Normalized()
	m := inventoryBaseQuery(ctx)
	if skuId > 0 {
		m = m.Where("i.sku_id", skuId)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		m = m.Where("(s.sku_no LIKE ? OR s.name LIKE ?)", like, like)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计库存失败")
	}
	recs, err := m.Fields(inventoryCols).Order("i.sku_id ASC").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询库存失败")
	}
	list := make([]model.InventoryItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, inventoryFromRecord(r))
	}
	return &model.PageResult[model.InventoryItem]{List: list, Total: int64(total)}, nil
}

// InventoryWarnings 预警列表（FR-007）: 可售 <= 预警阈值（含等于）。
func InventoryWarnings(ctx context.Context, page model.PageReq) (*model.PageResult[model.InventoryItem], error) {
	page = page.Normalized()
	m := inventoryBaseQuery(ctx).Where("i.total - i.locked <= i.warn_count")
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计预警失败")
	}
	recs, err := m.Fields(inventoryCols).Order("i.total - i.locked ASC, i.sku_id ASC").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询预警失败")
	}
	list := make([]model.InventoryItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, inventoryFromRecord(r))
	}
	return &model.PageResult[model.InventoryItem]{List: list, Total: int64(total)}, nil
}

// InventoryAdjust 库存调整（FR-008）: delta 有符号; 不可为负（条件更新）; 同事务留痕; 返回调整后总量。
func InventoryAdjust(ctx context.Context, skuId int64, delta int, operator, remark string) (int, error) {
	if delta == 0 {
		return 0, errcode.New(errcode.CodeInvalidParam, "调整数量不能为 0")
	}
	cols := dao.Inventory.Columns()
	var after int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 条件更新（防负）: inventory.total 为 INT UNSIGNED——不能写 `total + (-n) >= 0`
		// （无符号负数运算触发 out of range, 错误码 52）。按 delta 符号分支:
		//   delta >= 0: 无下溢风险, 仅判存在;
		//   delta < 0 : 判「可售 >= 扣减量」（010 评审 I1: 仅判 total>=扣减量会漏护可售,
		//               撞到表级 CHECK(chk_locked_le_total) 变成无业务码的系统错误）。
		q := dao.Inventory.Ctx(ctx).Where(cols.SkuId, skuId)
		if delta < 0 {
			q = q.Where("total - locked >= ?", -delta)
		}
		res, e := q.
			Data(do.Inventory{Total: gdb.Raw(fmt.Sprintf("total + (%d)", delta))}).
			Update()
		if e != nil {
			return gerror.Wrap(e, "调整库存失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeInventoryAdjust, "库存调整非法: SKU不存在或调整后为负")
		}
		// 事务内读回快照（对账用）
		rec, e := dao.Inventory.Ctx(ctx).
			Fields(cols.Total, cols.Locked).
			Where(cols.SkuId, skuId).One()
		if e != nil {
			return gerror.Wrap(e, "读取库存失败")
		}
		if rec.IsEmpty() {
			return errcode.New(errcode.CodeInventoryAdjust, "库存记录不存在")
		}
		after = rec[cols.Total].Int()
		lockedAfter := rec[cols.Locked].Int()

		qty := delta
		if qty < 0 {
			qty = -qty
		}
		if _, e = dao.InventoryLog.Ctx(ctx).Data(do.InventoryLog{
			SkuId:       skuId,
			ChangeType:  5, // 5=后台调整
			Quantity:    qty,
			TotalAfter:  after,
			LockedAfter: lockedAfter,
			Operator:    operator,
			Remark:      remark,
		}).Insert(); e != nil {
			return gerror.Wrap(e, "写入库存流水失败")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return after, nil
}
