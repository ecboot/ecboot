package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminUserCreate 创建后台账号
func (c *ControllerV1) AdminUserCreate(ctx context.Context, req *v1.AdminUserCreateReq) (res *v1.AdminUserCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
