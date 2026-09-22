package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// DistAccount 佣金账户（可负余额）。
func (c *ControllerV1) DistAccount(ctx context.Context, req *v1.DistAccountReq) (res *v1.DistAccountRes, err error) {
	acc, err := user.NewDistributionLogic().Account(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.DistAccountRes{Balance: acc.Balance, Frozen: acc.Frozen}, nil
}
