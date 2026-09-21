package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// FavoriteAdd 收藏商品（已软删则复活; 幂等）
func (c *ControllerV1) FavoriteAdd(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	if err = user.FavoriteAdd(ctx, middleware.CtxUserIdFrom(ctx), spuId); err != nil {
		return nil, err
	}
	return &v1.FavoriteAddRes{Success: true}, nil
}
