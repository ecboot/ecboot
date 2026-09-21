package common

import (
	"context"

	"ecboot/api/common/v1"
	"ecboot/internal/library/sms"
)

func (c *ControllerV1) GetSmsCode(ctx context.Context, req *v1.GetSmsCodeReq) (res *v1.GetSmsCodeRes, err error) {
	expiresIn, err := sms.SendSmsCode(ctx, req.PhoneNumber, req.CaptchaKey, req.CaptchaCode)
	if err != nil {
		return nil, err
	}
	return &v1.GetSmsCodeRes{ExpiresIn: expiresIn}, nil
}
