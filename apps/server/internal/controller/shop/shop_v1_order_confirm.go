package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// OrderConfirm 确认收货（仅待收货）
func (c *ControllerV1) OrderConfirm(ctx context.Context, req *v1.OrderConfirmReq) (res *v1.OrderConfirmRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	if err = shop.NewOrderLogic().Confirm(ctx, userId, req.OrderNo); err != nil {
		return nil, err
	}
	return &v1.OrderConfirmRes{Success: true}, nil
}
