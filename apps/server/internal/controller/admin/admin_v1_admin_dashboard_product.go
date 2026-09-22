package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminDashboardProduct 商品看板（FR-3）: 在售/低库存预警/待审核评价。
func (c *ControllerV1) AdminDashboardProduct(ctx context.Context, req *v1.AdminDashboardProductReq) (res *v1.AdminDashboardProductRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDashboardRead); err != nil {
		return nil, err
	}
	out, err := system.NewDashboardLogic().Product(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.AdminDashboardProductRes{
		OnSaleCount: out.OnSaleCount, LowStockCount: out.LowStockCount, PendingReview: out.PendingReview,
	}, nil
}
