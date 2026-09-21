package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/shop"
)

// AdminOrderDetail 后台订单详情（含卖家备注——买家不可见）
func (c *ControllerV1) AdminOrderDetail(ctx context.Context, req *v1.AdminOrderDetailReq) (res *v1.AdminOrderDetailRes, err error) {
	d, err := shop.NewOrderLogic().AdminDetail(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	res = &v1.AdminOrderDetailRes{
		OrderNo: d.OrderNo, Status: d.Status, RefundStatus: d.RefundStatus,
		Amount: v1.AdminOrderAmountBrief{
			TotalAmount: d.Amount.TotalAmount, PromotionAmount: d.Amount.CouponAmount +
				d.Amount.FullReductionAmount + d.Amount.PointAmount,
			FreightAmount: d.Amount.FreightAmount, AccountAmount: d.Amount.AccountAmount,
			PayAmount: d.Amount.PayAmount,
		},
		Receiver: d.Receiver, UserRemark: d.UserRemark, SellerRemark: d.SellerRemark,
		Pay: d.Pay, Deliver: d.Deliver,
		Items:      make([]v1.AdminOrderItemBrief, 0, len(d.Items)),
		StatusLogs: make([]v1.AdminOrderStatusLog, 0, len(d.StatusLogs)),
	}
	for _, li := range d.Items {
		res.Items = append(res.Items, v1.AdminOrderItemBrief{
			SpuName: li.SpuName, SkuSpecs: li.SkuSpecs, Image: li.Image,
			Quantity: li.Quantity, Price: li.Price,
		})
	}
	for _, lg := range d.StatusLogs {
		res.StatusLogs = append(res.StatusLogs, v1.AdminOrderStatusLog{
			FromStatus: lg.FromStatus, ToStatus: lg.ToStatus,
			Remark: lg.Remark, CreatedAt: lg.CreatedAt,
		})
	}
	return res, nil
}
