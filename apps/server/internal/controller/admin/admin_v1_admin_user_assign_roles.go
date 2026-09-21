package admin

import (
	"context"
	"strconv"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/system"
)

// AdminUserAssignRoles 账号-角色分配
func (c *ControllerV1) AdminUserAssignRoles(ctx context.Context, req *v1.AdminUserAssignRolesReq) (res *v1.AdminUserAssignRolesRes, err error) {
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}
	roleIds := make([]int64, 0, len(req.RoleIds))
	for _, s := range req.RoleIds {
		rid, perr := strconv.ParseInt(s, 10, 64)
		if perr != nil {
			return nil, perr
		}
		roleIds = append(roleIds, rid)
	}
	if err = system.AssignRoles(ctx, id, roleIds); err != nil {
		return nil, err
	}
	return &v1.AdminUserAssignRolesRes{Success: true}, nil
}
