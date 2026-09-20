package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// ShareCode 我的推广码
func (c *ControllerV1) ShareCode(ctx context.Context, req *v1.ShareCodeReq) (res *v1.ShareCodeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
