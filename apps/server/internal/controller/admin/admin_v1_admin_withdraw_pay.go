package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminWithdrawPay 打款登记（成功 40 渠道单号幂等 / 失败 60 回退）。
func (c *ControllerV1) AdminWithdrawPay(ctx context.Context, req *v1.AdminWithdrawPayReq) (res *v1.AdminWithdrawPayRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionWithdrawPay); err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminWithdrawPay(ctx, req.WithdrawNo, req.Success, req.ChannelOrderNo, req.FailReason); err != nil {
		return nil, err
	}
	return &v1.AdminWithdrawPayRes{Success: true}, nil
}
