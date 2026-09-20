package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

func (c *ControllerV1) FavoriteList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
