package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRiskRuleCreate 新增风控规则。
func (c *ControllerV1) AdminRiskRuleCreate(ctx context.Context, req *v1.AdminRiskRuleCreateReq) (res *v1.AdminRiskRuleCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRuleManage); err != nil {
		return nil, err
	}
	id, err := system.NewRiskAdminLogic().AdminRuleCreate(ctx, req.Name, req.RuleType, req.ConditionExpr, req.Action)
	if err != nil {
		return nil, err
	}
	return &v1.AdminRiskRuleCreateRes{Id: fmtID(id)}, nil
}
