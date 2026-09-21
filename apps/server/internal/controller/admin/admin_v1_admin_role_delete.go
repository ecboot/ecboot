package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRoleDelete 删除角色(软删)
func (c *ControllerV1) AdminRoleDelete(ctx context.Context, req *v1.AdminRoleDeleteReq) (res *v1.AdminRoleDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:role:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.RoleDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminRoleDeleteRes{Success: true}, nil
}
