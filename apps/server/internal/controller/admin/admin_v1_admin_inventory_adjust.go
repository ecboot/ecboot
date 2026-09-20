package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminInventoryAdjust 库存调整（留痕 inventory_log）
func (c *ControllerV1) AdminInventoryAdjust(ctx context.Context, req *v1.AdminInventoryAdjustReq) (res *v1.AdminInventoryAdjustRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
