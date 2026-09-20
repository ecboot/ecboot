package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAfterSaleReject 拒绝
func (c *ControllerV1) AdminAfterSaleReject(ctx context.Context, req *v1.AdminAfterSaleRejectReq) (res *v1.AdminAfterSaleRejectRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
