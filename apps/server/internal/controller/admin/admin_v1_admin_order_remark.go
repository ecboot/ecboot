package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminOrderRemark 内部备注（买家不可见）
func (c *ControllerV1) AdminOrderRemark(ctx context.Context, req *v1.AdminOrderRemarkReq) (res *v1.AdminOrderRemarkRes, err error) {
	if err = middleware.RequirePerm(ctx, "order:update"); err != nil {
		return nil, err
	}
	operator := adminOperator(ctx)
	if err = shop.NewOrderLogic().SellerRemark(ctx, req.OrderNo, req.SellerRemark, operator); err != nil {
		return nil, err
	}
	return &v1.AdminOrderRemarkRes{Success: true}, nil
}
