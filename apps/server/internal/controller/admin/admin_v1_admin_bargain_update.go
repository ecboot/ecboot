package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminBargainUpdate(ctx context.Context, req *v1.AdminBargainUpdateReq) (res *v1.AdminBargainUpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
