package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminProfile 个人信息
func (c *ControllerV1) AdminProfile(ctx context.Context, req *v1.AdminProfileReq) (res *v1.AdminProfileRes, err error) {
	p, err := system.Profile(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.AdminProfileRes{
		Username: p.Username,
		RealName: p.RealName,
		Roles:    p.Roles,
	}, nil
}
