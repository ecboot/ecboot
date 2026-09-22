package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminDistRuleDelete 软删佣金规则。
func (c *ControllerV1) AdminDistRuleDelete(ctx context.Context, req *v1.AdminDistRuleDeleteReq) (res *v1.AdminDistRuleDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRuleManage); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.NewDistributionAdminLogic().AdminRuleDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminDistRuleDeleteRes{Success: true}, nil
}
