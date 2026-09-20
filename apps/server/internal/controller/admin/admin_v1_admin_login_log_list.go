package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminLoginLogList 后台登录审计
func (c *ControllerV1) AdminLoginLogList(ctx context.Context, req *v1.AdminLoginLogListReq) (res *v1.AdminLoginLogListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
