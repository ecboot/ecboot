package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminCouponList ---------- 优惠券模板 ----------
func (c *ControllerV1) AdminCouponList(ctx context.Context, req *v1.AdminCouponListReq) (res *v1.AdminCouponListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
