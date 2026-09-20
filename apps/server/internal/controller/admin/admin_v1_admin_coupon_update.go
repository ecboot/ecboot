package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminCouponUpdate(ctx context.Context, req *v1.AdminCouponUpdateReq) (res *v1.AdminCouponUpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
