package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// DistStatus 推广员状态与等级
func (c *ControllerV1) DistStatus(ctx context.Context, req *v1.DistStatusReq) (res *v1.DistStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
