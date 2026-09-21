package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
	"ecboot/internal/library/security"
)

func (c *ControllerV1) AdminTokenRefresh(ctx context.Context, req *v1.AdminTokenRefreshReq) (res *v1.AdminTokenRefreshRes, err error) {
	sm := security.NewSessionManager("admin", 7)
	token, refresh, _, err := sm.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	return &v1.AdminTokenRefreshRes{Token: token, RefreshToken: refresh}, nil
}
