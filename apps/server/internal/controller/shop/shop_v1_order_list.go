package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// OrderList 我的订单列表
func (c *ControllerV1) OrderList(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewOrderLogic().List(ctx, userId, req.Status,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.OrderListRes{List: make([]v1.OrderListItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		item := v1.OrderListItem{
			OrderNo: it.OrderNo, Status: it.Status, CreatedAt: it.CreatedAt,
			Amount: v1.OrderAmountBrief{
				TotalAmount: it.Amount.TotalAmount, PromotionAmount: it.Amount.CouponAmount +
					it.Amount.FullReductionAmount,
				FreightAmount: it.Amount.FreightAmount, AccountAmount: it.Amount.AccountAmount,
				PayAmount: it.Amount.PayAmount,
			},
			Items: make([]v1.OrderItemBrief, 0, len(it.Items)),
		}
		for _, li := range it.Items {
			item.Items = append(item.Items, v1.OrderItemBrief{
				SpuName: li.SpuName, SkuSpecs: li.SkuSpecs, Image: li.Image,
				Quantity: li.Quantity, Price: li.Price,
			})
		}
		res.List = append(res.List, item)
	}
	return res, nil
}
