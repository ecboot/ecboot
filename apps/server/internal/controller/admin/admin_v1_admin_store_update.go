package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminStoreUpdate(ctx context.Context, req *v1.AdminStoreUpdateReq) (res *v1.AdminStoreUpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
