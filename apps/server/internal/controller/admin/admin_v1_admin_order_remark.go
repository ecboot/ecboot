package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOrderRemark 卖家备注
func (c *ControllerV1) AdminOrderRemark(ctx context.Context, req *v1.AdminOrderRemarkReq) (res *v1.AdminOrderRemarkRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
