package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminCouponRecordList 券发放/使用记录
func (c *ControllerV1) AdminCouponRecordList(ctx context.Context, req *v1.AdminCouponRecordListReq) (res *v1.AdminCouponRecordListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
