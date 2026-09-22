package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// ShareCode 我的推广码（稳定）。
func (c *ControllerV1) ShareCode(ctx context.Context, req *v1.ShareCodeReq) (res *v1.ShareCodeRes, err error) {
	code, link, err := user.NewDistributionLogic().ShareCode(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.ShareCodeRes{ShareCode: code, ShareLink: link}, nil
}
