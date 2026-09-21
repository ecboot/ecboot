package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// PointAccount 积分账户（余额可负, 如实展示）
func (c *ControllerV1) PointAccount(ctx context.Context, req *v1.PointAccountReq) (res *v1.PointAccountRes, err error) {
	acc, err := user.PointAccount(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.PointAccountRes{Balance: acc.Balance, LastEarnedAt: acc.LastEarnedAt}, nil
}
