package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminRoleCreate 新增角色
func (c *ControllerV1) AdminRoleCreate(ctx context.Context, req *v1.AdminRoleCreateReq) (res *v1.AdminRoleCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:role:manage"); err != nil {
		return nil, err
	}
	id, err := system.RoleCreate(ctx, model.RoleInput{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminRoleCreateRes{Id: fmtID(id)}, nil
}
