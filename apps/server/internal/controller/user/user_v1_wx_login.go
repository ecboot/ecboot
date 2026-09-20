package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// WxLogin 微信登录（手机号优先归并; 开发态 mock）
func (c *ControllerV1) WxLogin(ctx context.Context, req *v1.WxLoginReq) (res *v1.WxLoginRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
