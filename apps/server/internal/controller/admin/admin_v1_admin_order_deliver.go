package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminOrderDeliver 发货（物流公司须存在且启用）
func (c *ControllerV1) AdminOrderDeliver(ctx context.Context, req *v1.AdminOrderDeliverReq) (res *v1.AdminOrderDeliverRes, err error) {
	if err = middleware.RequirePerm(ctx, "order:deliver"); err != nil {
		return nil, err
	}
	operator := adminOperator(ctx)
	if err = shop.NewOrderLogic().Deliver(ctx, req.OrderNo, req.LogisticsCode, req.DeliverNo, operator); err != nil {
		return nil, err
	}
	return &v1.AdminOrderDeliverRes{Success: true}, nil
}
