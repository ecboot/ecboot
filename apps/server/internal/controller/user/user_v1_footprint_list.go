package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

func (c *ControllerV1) FootprintList(ctx context.Context, req *v1.FootprintListReq) (res *v1.FootprintListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
