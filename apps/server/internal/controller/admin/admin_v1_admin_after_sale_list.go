package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleList 售后列表（状态筛选）
// 权限: aftersale:read
func (c *ControllerV1) AdminAfterSaleList(ctx context.Context, req *v1.AdminAfterSaleListReq) (res *v1.AdminAfterSaleListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleRead); err != nil {
		return nil, err
	}
	out, err := shop.NewAfterSaleLogic().AdminList(ctx, req.Status, req.Type,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminAfterSaleListRes{List: make([]v1.AdminAfterSaleItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminAfterSaleItem{
			AfterSaleNo:  it.AfterSaleNo,
			OrderNo:      it.OrderNo,
			UserId:       fmtID(it.UserId),
			Type:         it.Type,
			Quantity:     it.Quantity,
			RefundAmount: it.RefundAmount,
			Status:       it.Status,
			CreatedAt:    it.CreatedAt,
		})
	}
	return res, nil
}
