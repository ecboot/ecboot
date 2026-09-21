package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// FootprintClear 清空足迹
func (c *ControllerV1) FootprintClear(ctx context.Context, req *v1.FootprintClearReq) (res *v1.FootprintClearRes, err error) {
	if err = user.FootprintClear(ctx, middleware.CtxUserIdFrom(ctx)); err != nil {
		return nil, err
	}
	return &v1.FootprintClearRes{Success: true}, nil
}
