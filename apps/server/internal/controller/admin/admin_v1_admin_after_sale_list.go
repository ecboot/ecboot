package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAfterSaleList 售后列表（状态筛选）
func (c *ControllerV1) AdminAfterSaleList(ctx context.Context, req *v1.AdminAfterSaleListReq) (res *v1.AdminAfterSaleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
