package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRiskAppeal 申诉处理（0→2通过 / 0→3驳回; 已结论拒绝）。
func (c *ControllerV1) AdminRiskAppeal(ctx context.Context, req *v1.AdminRiskAppealReq) (res *v1.AdminRiskAppealRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRecordAppeal); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.NewRiskAdminLogic().AdminAppeal(ctx, id, req.Pass, req.Remark); err != nil {
		return nil, err
	}
	return &v1.AdminRiskAppealRes{Success: true}, nil
}
