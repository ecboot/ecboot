package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// DistStatus 我的推广员状态。
func (c *ControllerV1) DistStatus(ctx context.Context, req *v1.DistStatusReq) (res *v1.DistStatusRes, err error) {
	st, err := user.NewDistributionLogic().Status(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.DistStatusRes{Status: st.Status, Level: st.Level}, nil
}
