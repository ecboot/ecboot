package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminFlashSaleItems(ctx context.Context, req *v1.AdminFlashSaleItemsReq) (res *v1.AdminFlashSaleItemsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
