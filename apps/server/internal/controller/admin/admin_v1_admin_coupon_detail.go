package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminCouponDetail 券模板详情
func (c *ControllerV1) AdminCouponDetail(ctx context.Context, req *v1.AdminCouponDetailReq) (res *v1.AdminCouponDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
