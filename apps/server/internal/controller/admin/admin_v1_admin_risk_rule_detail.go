package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRiskRuleDetail 风控规则详情。
func (c *ControllerV1) AdminRiskRuleDetail(ctx context.Context, req *v1.AdminRiskRuleDetailReq) (res *v1.AdminRiskRuleDetailRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRuleRead); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	it, err := system.NewRiskAdminLogic().AdminRuleDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminRiskRuleDetailRes{
		Id: fmtID(it.Id), Name: it.Name, RuleType: it.RuleType,
		ConditionExpr: it.ConditionExpr, Action: it.Action, Status: it.Status,
	}, nil
}
