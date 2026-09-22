package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// DistApply 申请成为推广员（重复申请 60001）。
func (c *ControllerV1) DistApply(ctx context.Context, req *v1.DistApplyReq) (res *v1.DistApplyRes, err error) {
	if err = user.NewDistributionLogic().Apply(ctx, middleware.CtxUserIdFrom(ctx)); err != nil {
		return nil, err
	}
	return &v1.DistApplyRes{Status: 1}, nil
}
