package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminDistRuleUpdate 修改佣金规则（status 三态: 不传不修改）。
func (c *ControllerV1) AdminDistRuleUpdate(ctx context.Context, req *v1.AdminDistRuleUpdateReq) (res *v1.AdminDistRuleUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRuleManage); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminRuleUpdate(ctx, id, req.Level1Rate, req.Level2Rate, req.Status); err != nil {
		return nil, err
	}
	return &v1.AdminDistRuleUpdateRes{Success: true}, nil
}
