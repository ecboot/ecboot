package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminConfigUpdate 修改配置
func (c *ControllerV1) AdminConfigUpdate(ctx context.Context, req *v1.AdminConfigUpdateReq) (res *v1.AdminConfigUpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
