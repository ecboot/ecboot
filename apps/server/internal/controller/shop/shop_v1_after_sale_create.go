package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AfterSaleCreate 申请售后（按订单项; 仅退款/退货退款）
func (c *ControllerV1) AfterSaleCreate(ctx context.Context, req *v1.AfterSaleCreateReq) (res *v1.AfterSaleCreateRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	itemId, err := parseID(req.OrderItemId)
	if err != nil {
		return nil, err
	}
	no, err := shop.NewAfterSaleLogic().Apply(ctx, userId, model.AfterSaleApplyInput{
		OrderItemId:   itemId,
		Type:          req.Type,
		Quantity:      req.Quantity,
		Reason:        req.Reason,
		Description:   req.Description,
		VoucherImages: req.VoucherImages,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AfterSaleCreateRes{AfterSaleNo: no}, nil
}
