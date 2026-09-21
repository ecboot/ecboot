package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// OrderCancel 取消订单（仅待付款; 释放库存/券/余额）
func (c *ControllerV1) OrderCancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	if err = shop.NewOrderLogic().Cancel(ctx, userId, req.OrderNo, req.Reason); err != nil {
		return nil, err
	}
	return &v1.OrderCancelRes{Success: true}, nil
}
