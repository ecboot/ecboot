package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOrderDeliver 发货
func (c *ControllerV1) AdminOrderDeliver(ctx context.Context, req *v1.AdminOrderDeliverReq) (res *v1.AdminOrderDeliverRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
