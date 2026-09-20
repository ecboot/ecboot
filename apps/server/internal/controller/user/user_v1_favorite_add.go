package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// FavoriteAdd 收藏（已软删则复活, 重复收藏幂等）
func (c *ControllerV1) FavoriteAdd(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
