package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// CouponAvailableList 可领模板列表
func (c *ControllerV1) CouponAvailableList(ctx context.Context, req *v1.CouponAvailableListReq) (res *v1.CouponAvailableListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
