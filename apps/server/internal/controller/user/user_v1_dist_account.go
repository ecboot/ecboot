package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// DistAccount 佣金账户
func (c *ControllerV1) DistAccount(ctx context.Context, req *v1.DistAccountReq) (res *v1.DistAccountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
