package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

func (c *ControllerV1) GetSmsCode(ctx context.Context, req *v1.GetSmsCodeReq) (res *v1.GetSmsCodeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
