package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminDashboardTrade 交易看板（FR-1）: 订单数/销售额/退款额/待发货。
func (c *ControllerV1) AdminDashboardTrade(ctx context.Context, req *v1.AdminDashboardTradeReq) (res *v1.AdminDashboardTradeRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDashboardRead); err != nil {
		return nil, err
	}
	out, err := system.NewDashboardLogic().Trade(ctx, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	return &v1.AdminDashboardTradeRes{
		OrderCount: out.OrderCount, SalesAmount: out.SalesAmount,
		RefundAmount: out.RefundAmount, PendingDeliver: out.PendingDeliver,
	}, nil
}
