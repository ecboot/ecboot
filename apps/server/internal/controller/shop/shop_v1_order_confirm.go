package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// OrderConfirm 确认收货（仅待收货）
func (c *ControllerV1) OrderConfirm(ctx context.Context, req *v1.OrderConfirmReq) (res *v1.OrderConfirmRes, err error) {
	if err = shop.NewOrderLogic().Confirm(ctx, middleware.CtxUserIdFrom(ctx), req.OrderNo); err != nil {
		return nil, err
	}
	return &v1.OrderConfirmRes{Success: true}, nil
}
