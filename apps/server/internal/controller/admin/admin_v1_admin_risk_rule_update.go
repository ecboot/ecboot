package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRiskRuleUpdate 修改风控规则（status 三态: 不传不修改）。
func (c *ControllerV1) AdminRiskRuleUpdate(ctx context.Context, req *v1.AdminRiskRuleUpdateReq) (res *v1.AdminRiskRuleUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRuleManage); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.NewRiskAdminLogic().AdminRuleUpdate(ctx, id, req.Name, req.ConditionExpr, req.Action, req.Status); err != nil {
		return nil, err
	}
	return &v1.AdminRiskRuleUpdateRes{Success: true}, nil
}
