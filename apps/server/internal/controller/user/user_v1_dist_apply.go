package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// DistApply 申请成为推广员
func (c *ControllerV1) DistApply(ctx context.Context, req *v1.DistApplyReq) (res *v1.DistApplyRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
