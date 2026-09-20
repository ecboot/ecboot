package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOrderList 订单列表（多条件筛选）
func (c *ControllerV1) AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (res *v1.AdminOrderListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
