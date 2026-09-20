package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminBargainList ---------- 砍价 ----------
func (c *ControllerV1) AdminBargainList(ctx context.Context, req *v1.AdminBargainListReq) (res *v1.AdminBargainListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
