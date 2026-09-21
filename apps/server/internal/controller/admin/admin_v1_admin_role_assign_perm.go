package admin

import (
	"context"
	"strconv"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRoleAssignPerm 角色权限分配(全量替换)
func (c *ControllerV1) AdminRoleAssignPerm(ctx context.Context, req *v1.AdminRoleAssignPermReq) (res *v1.AdminRoleAssignPermRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:role:assign"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	permIds := make([]int64, 0, len(req.PermissionIds))
	for _, s := range req.PermissionIds {
		pid, perr := strconv.ParseInt(s, 10, 64)
		if perr != nil {
			return nil, perr
		}
		permIds = append(permIds, pid)
	}
	if err = system.AssignPermissions(ctx, id, permIds); err != nil {
		return nil, err
	}
	return &v1.AdminRoleAssignPermRes{Success: true}, nil
}
