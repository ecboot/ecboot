package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminAfterSaleDetail(ctx context.Context, req *v1.AdminAfterSaleDetailReq) (res *v1.AdminAfterSaleDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
