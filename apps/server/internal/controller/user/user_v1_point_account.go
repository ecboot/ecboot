package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

func (c *ControllerV1) PointAccount(ctx context.Context, req *v1.PointAccountReq) (res *v1.PointAccountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
