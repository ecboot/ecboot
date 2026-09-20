package user

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/api/user/v1"
	"ecboot/internal/service/user"
)

func (c *ControllerV1) SmsLogin(ctx context.Context, req *v1.SmsLoginReq) (res *v1.SmsLoginRes, err error) {
	out, err := user.SmsLogin(ctx, req.PhoneNumber, req.SmsCode, req.Channel)
	if err != nil {
		return nil, err
	}
	return &v1.SmsLoginRes{Token: out.Token, RefreshToken: out.RefreshToken, UserId: gconv.String(out.UserId), IsNew: out.IsNew}, nil
}
