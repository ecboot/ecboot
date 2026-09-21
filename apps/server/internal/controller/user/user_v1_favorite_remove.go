package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// FavoriteRemove 取消收藏（软删; 幂等）
func (c *ControllerV1) FavoriteRemove(ctx context.Context, req *v1.FavoriteRemoveReq) (res *v1.FavoriteRemoveRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	if err = user.FavoriteRemove(ctx, middleware.CtxUserIdFrom(ctx), spuId); err != nil {
		return nil, err
	}
	return &v1.FavoriteRemoveRes{Success: true}, nil
}
