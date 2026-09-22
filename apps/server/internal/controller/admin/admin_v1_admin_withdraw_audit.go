package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminWithdrawAudit 提现审核（通过 20 / 拒绝 50 回退; 条件状态机）。
func (c *ControllerV1) AdminWithdrawAudit(ctx context.Context, req *v1.AdminWithdrawAuditReq) (res *v1.AdminWithdrawAuditRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionWithdrawAudit); err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminWithdrawAudit(ctx, req.WithdrawNo, req.Pass, req.Reason); err != nil {
		return nil, err
	}
	return &v1.AdminWithdrawAuditRes{Success: true}, nil
}
