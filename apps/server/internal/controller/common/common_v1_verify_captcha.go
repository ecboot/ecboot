package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
)

func (c *ControllerV1) VerifyCaptcha(ctx context.Context, req *v1.VerifyCaptchaReq) (res *v1.VerifyCaptchaRes, err error) {
	ok, err := captcha.Verify(ctx, req.CaptchaKey, req.CaptchaCode)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, gerror.NewCode(gcode.New(errcode.CodeCaptchaError, "验证码错误或已过期", nil))
	}
	return &v1.VerifyCaptchaRes{Success: true}, nil
}
