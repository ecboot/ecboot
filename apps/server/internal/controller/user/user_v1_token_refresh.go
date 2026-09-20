package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// TokenRefresh 刷新访问凭证（双凭证会话）
func (c *ControllerV1) TokenRefresh(ctx context.Context, req *v1.TokenRefreshReq) (res *v1.TokenRefreshRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
