package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

func (c *ControllerV1) MyCouponList(ctx context.Context, req *v1.MyCouponListReq) (res *v1.MyCouponListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
