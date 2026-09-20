package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminUserAssignRoles 账号-角色分配
func (c *ControllerV1) AdminUserAssignRoles(ctx context.Context, req *v1.AdminUserAssignRolesReq) (res *v1.AdminUserAssignRolesRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
