package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
	"ecboot/internal/library/security"
	"ecboot/internal/service/system"
)

func (c *ControllerV1) AdminTokenRefresh(ctx context.Context, req *v1.AdminTokenRefreshReq) (res *v1.AdminTokenRefreshRes, err error) {
	sm := security.NewSessionManager("admin", 7)
	token, refresh, userId, err := sm.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	// 评审 M4: 被禁用/软删账号不得续新凭证对
	if ok, aErr := system.ActiveAdmin(ctx, userId); aErr != nil || !ok {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	return &v1.AdminTokenRefreshRes{Token: token, RefreshToken: refresh}, nil
}
