package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminConfigUpdate 修改系统配置
func (c *ControllerV1) AdminConfigUpdate(ctx context.Context, req *v1.AdminConfigUpdateReq) (res *v1.AdminConfigUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "system:config:update"); err != nil {
		return nil, err
	}
	if err = system.ConfigUpdate(ctx, req.Code, req.Value, req.Status); err != nil {
		return nil, err
	}
	return &v1.AdminConfigUpdateRes{Success: true}, nil
}
