package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// WithdrawCreate 提现申请（余额条件冻结; 并发只成一笔）。
func (c *ControllerV1) WithdrawCreate(ctx context.Context, req *v1.WithdrawCreateReq) (res *v1.WithdrawCreateRes, err error) {
	no, err := user.NewDistributionLogic().WithdrawApply(ctx, middleware.CtxUserIdFrom(ctx), req.Amount)
	if err != nil {
		return nil, err
	}
	return &v1.WithdrawCreateRes{WithdrawNo: no}, nil
}
