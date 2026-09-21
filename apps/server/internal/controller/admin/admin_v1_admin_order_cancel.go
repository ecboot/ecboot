package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminOrderCancel 取消订单（与 C 端同语义: 状态 90 + 资源释放）
func (c *ControllerV1) AdminOrderCancel(ctx context.Context, req *v1.AdminOrderCancelReq) (res *v1.AdminOrderCancelRes, err error) {
	if err = middleware.RequirePerm(ctx, "order:cancel"); err != nil {
		return nil, err
	}
	operator := adminOperator(ctx)
	if err = shop.NewOrderLogic().AdminCancel(ctx, req.OrderNo, req.Reason, operator); err != nil {
		return nil, err
	}
	return &v1.AdminOrderCancelRes{Success: true}, nil
}
