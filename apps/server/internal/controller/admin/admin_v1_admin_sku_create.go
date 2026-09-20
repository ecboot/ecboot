package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminSkuCreate 新增 SKU
func (c *ControllerV1) AdminSkuCreate(ctx context.Context, req *v1.AdminSkuCreateReq) (res *v1.AdminSkuCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
