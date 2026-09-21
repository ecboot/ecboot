package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminChangePassword 修改密码
func (c *ControllerV1) AdminChangePassword(ctx context.Context, req *v1.AdminChangePasswordReq) (res *v1.AdminChangePasswordRes, err error) {
	if err = system.ChangePassword(ctx, middleware.CtxUserIdFrom(ctx), req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}
	return &v1.AdminChangePasswordRes{Success: true}, nil
}
