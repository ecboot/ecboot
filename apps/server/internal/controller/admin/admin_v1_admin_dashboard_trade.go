package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDashboardTrade 交易看板
func (c *ControllerV1) AdminDashboardTrade(ctx context.Context, req *v1.AdminDashboardTradeReq) (res *v1.AdminDashboardTradeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
