package user

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/api/user/v1"
	"ecboot/internal/service/user"
)

// WxLogin 微信登录（手机号优先归并; 开发态 mock）
func (c *ControllerV1) WxLogin(ctx context.Context, req *v1.WxLoginReq) (res *v1.WxLoginRes, err error) {
	out, err := user.WxLogin(ctx, req.WxCode, req.Phone, req.SmsCode, req.Channel)
	if err != nil {
		return nil, err
	}
	return &v1.WxLoginRes{Token: out.Token, RefreshToken: out.RefreshToken, UserId: gconv.String(out.UserId), IsNew: out.IsNew}, nil
}
