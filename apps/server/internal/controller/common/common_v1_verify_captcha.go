package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

// VerifyCaptcha 图形验证码独立校验（一次性; 随 003 实现补齐逻辑）
func (c *ControllerV1) VerifyCaptcha(ctx context.Context, req *v1.VerifyCaptchaReq) (res *v1.VerifyCaptchaRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
