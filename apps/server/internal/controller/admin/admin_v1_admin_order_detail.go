package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOrderDetail 订单详情（含收货快照/状态流水）
func (c *ControllerV1) AdminOrderDetail(ctx context.Context, req *v1.AdminOrderDetailReq) (res *v1.AdminOrderDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
