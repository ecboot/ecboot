package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOrderCancel 管理员取消（待付款）
func (c *ControllerV1) AdminOrderCancel(ctx context.Context, req *v1.AdminOrderCancelReq) (res *v1.AdminOrderCancelRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
