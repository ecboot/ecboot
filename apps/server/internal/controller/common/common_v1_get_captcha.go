package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

func (c *ControllerV1) GetCaptcha(ctx context.Context, req *v1.GetCaptchaReq) (res *v1.GetCaptchaRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
