package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// OrderDetail 订单详情（他人订单按不存在）
func (c *ControllerV1) OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	d, err := shop.NewOrderLogic().OrderDetail(ctx, userId, req.OrderNo)
	if err != nil {
		return nil, err
	}
	res = &v1.OrderDetailRes{
		OrderNo: d.OrderNo, Status: d.Status, RefundStatus: d.RefundStatus,
		Amount: v1.OrderAmountBrief{
			TotalAmount: d.Amount.TotalAmount, PromotionAmount: d.Amount.CouponAmount +
				d.Amount.FullReductionAmount,
			FreightAmount: d.Amount.FreightAmount, AccountAmount: d.Amount.AccountAmount,
			PayAmount: d.Amount.PayAmount,
		},
		Receiver: d.Receiver, CreatedAt: d.CreatedAt,
		Items:       make([]v1.OrderItemBrief, 0, len(d.Items)),
		DeliverInfo: d.Deliver, CancelInfo: d.Cancel,
	}
	if d.Pay != nil {
		res.PayTime = d.Pay["successTime"]
	}
	for _, li := range d.Items {
		res.Items = append(res.Items, v1.OrderItemBrief{
			SpuName: li.SpuName, SkuSpecs: li.SkuSpecs, Image: li.Image,
			Quantity: li.Quantity, Price: li.Price,
		})
	}
	for _, lg := range d.StatusLogs {
		res.StatusLogs = append(res.StatusLogs, v1.OrderStatusLog{
			FromStatus: lg.FromStatus, ToStatus: lg.ToStatus,
			Remark: lg.Remark, CreatedAt: lg.CreatedAt,
		})
	}
	return res, nil
}
