package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/system"
)

// AdminRiskRuleDelete 软删风控规则。
func (c *ControllerV1) AdminRiskRuleDelete(ctx context.Context, req *v1.AdminRiskRuleDeleteReq) (res *v1.AdminRiskRuleDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRuleManage); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = system.NewRiskAdminLogic().AdminRuleDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminRiskRuleDeleteRes{Success: true}, nil
}
