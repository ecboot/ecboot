package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDashboardProduct 商品看板
func (c *ControllerV1) AdminDashboardProduct(ctx context.Context, req *v1.AdminDashboardProductReq) (res *v1.AdminDashboardProductRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
