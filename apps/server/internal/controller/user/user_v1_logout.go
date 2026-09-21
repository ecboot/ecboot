package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/library/security"
	"ecboot/internal/middleware"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	token := ctx.Value(middleware.CtxToken)
	if s, ok := token.(string); ok && s != "" {
		sm := security.NewSessionManager("user", 7)
		if err = sm.Destroy(ctx, s, req.RefreshToken); err != nil {
			return nil, err
		}
	}
	return &v1.LogoutRes{Success: true}, nil
}
