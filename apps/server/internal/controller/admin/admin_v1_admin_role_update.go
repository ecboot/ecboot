package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminRoleUpdate 修改角色
func (c *ControllerV1) AdminRoleUpdate(ctx context.Context, req *v1.AdminRoleUpdateReq) (res *v1.AdminRoleUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:role:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.RoleUpdate(ctx, id, model.RoleInput{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminRoleUpdateRes{Success: true}, nil
}
