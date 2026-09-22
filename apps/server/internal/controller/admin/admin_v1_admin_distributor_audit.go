package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminDistributorAudit 推广员审核。
func (c *ControllerV1) AdminDistributorAudit(ctx context.Context, req *v1.AdminDistributorAuditReq) (res *v1.AdminDistributorAuditRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionAudit); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminDistributorAudit(ctx, id, req.Pass); err != nil {
		return nil, err
	}
	return &v1.AdminDistributorAuditRes{Success: true}, nil
}
