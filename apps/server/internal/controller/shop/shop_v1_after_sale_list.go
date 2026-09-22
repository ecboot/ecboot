package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AfterSaleList 售后列表（状态筛选 + 分页; 仅本人）
func (c *ControllerV1) AfterSaleList(ctx context.Context, req *v1.AfterSaleListReq) (res *v1.AfterSaleListRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewAfterSaleLogic().List(ctx, userId, req.Status,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AfterSaleListRes{List: make([]v1.AfterSaleListItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AfterSaleListItem{
			AfterSaleNo:  it.AfterSaleNo,
			OrderNo:      it.OrderNo,
			Type:         it.Type,
			RefundAmount: it.RefundAmount,
			Status:       it.Status,
			CreatedAt:    it.CreatedAt,
		})
	}
	return res, nil
}
