package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminInventoryWarnings 库存预警列表（可售≤阈值）
func (c *ControllerV1) AdminInventoryWarn(ctx context.Context, req *v1.AdminInventoryWarnReq) (res *v1.AdminInventoryWarnRes, err error) {
	out, err := shop.InventoryWarnings(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminInventoryWarnRes{List: make([]v1.AdminInventoryItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminInventoryItem{
			SkuId: fmtID(it.SkuId), SkuNo: it.SkuNo, SkuName: it.SkuName,
			Total: it.Total, Locked: it.Locked, Available: it.Available, WarnCount: it.WarnCount,
		})
	}
	return res, nil
}
