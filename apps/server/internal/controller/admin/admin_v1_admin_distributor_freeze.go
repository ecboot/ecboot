package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminDistributorFreeze 冻结/解冻推广员。
func (c *ControllerV1) AdminDistributorFreeze(ctx context.Context, req *v1.AdminDistributorFreezeReq) (res *v1.AdminDistributorFreezeRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionAudit); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminDistributorFreeze(ctx, id, req.Freeze); err != nil {
		return nil, err
	}
	return &v1.AdminDistributorFreezeRes{Success: true}, nil
}
