package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
	"ecboot/internal/library/security"
)

// TokenRefresh 刷新访问凭证（双凭证会话）
func (c *ControllerV1) TokenRefresh(ctx context.Context, req *v1.TokenRefreshReq) (res *v1.TokenRefreshRes, err error) {
	sm := security.NewSessionManager(7)
	token, refresh, _, err := sm.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	return &v1.TokenRefreshRes{Token: token, RefreshToken: refresh}, nil
}
