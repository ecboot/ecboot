package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminUserDelete 删除后台账号(软删)
func (c *ControllerV1) AdminUserDelete(ctx context.Context, req *v1.AdminUserDeleteReq) (res *v1.AdminUserDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:admin:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.AdminUserDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminUserDeleteRes{Success: true}, nil
}
