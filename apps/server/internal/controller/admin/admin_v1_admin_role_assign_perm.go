package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRoleAssignPerm 角色-权限全量替换
func (c *ControllerV1) AdminRoleAssignPerm(ctx context.Context, req *v1.AdminRoleAssignPermReq) (res *v1.AdminRoleAssignPermRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
