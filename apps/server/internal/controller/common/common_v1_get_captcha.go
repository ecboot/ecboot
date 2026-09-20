package common

import (
	"context"

	"ecboot/internal/library/captcha"

	"ecboot/api/common/v1"
)

func (c *ControllerV1) GetCaptcha(ctx context.Context, req *v1.GetCaptchaReq) (res *v1.GetCaptchaRes, err error) {
	key, imgBase64, _, err := captcha.Generate(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCaptchaRes{
		CaptchaKey: key,
		CaptchaImg: imgBase64,
	}, nil
}
