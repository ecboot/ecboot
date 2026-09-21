package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminInventoryList 库存列表（可售=total-locked 推导）
func (c *ControllerV1) AdminInventoryList(ctx context.Context, req *v1.AdminInventoryListReq) (res *v1.AdminInventoryListRes, err error) {
	var skuId int64
	if req.SkuId != "" {
		if skuId, err = parseID(req.SkuId); err != nil {
			return nil, err
		}
	}
	out, err := shop.InventoryList(ctx, skuId, req.Keyword,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminInventoryListRes{List: make([]v1.AdminInventoryItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminInventoryItem{
			SkuId: fmtID(it.SkuId), SkuNo: it.SkuNo, SkuName: it.SkuName,
			Total: it.Total, Locked: it.Locked, Available: it.Available, WarnCount: it.WarnCount,
		})
	}
	return res, nil
}
