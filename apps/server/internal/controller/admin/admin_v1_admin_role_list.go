package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRoleList 角色 CRUD
func (c *ControllerV1) AdminRoleList(ctx context.Context, req *v1.AdminRoleListReq) (res *v1.AdminRoleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
