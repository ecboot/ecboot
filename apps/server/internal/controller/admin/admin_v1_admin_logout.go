package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/library/security"
	"ecboot/internal/middleware"
)

// AdminLogout 后台登出（双凭证同毁, FR-004）
func (c *ControllerV1) AdminLogout(ctx context.Context, req *v1.AdminLogoutReq) (res *v1.AdminLogoutRes, err error) {
	token := ctx.Value(middleware.CtxToken)
	if s, ok := token.(string); ok && s != "" {
		sm := security.NewSessionManager("admin", 7)
		if err = sm.Destroy(ctx, s, req.RefreshToken); err != nil {
			return nil, err
		}
	}
	return &v1.AdminLogoutRes{Success: true}, nil
}
