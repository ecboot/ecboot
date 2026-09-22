package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AdminDistRuleCreate 创建佣金规则（作用域唯一, 撞键转业务码）。
func (c *ControllerV1) AdminDistRuleCreate(ctx context.Context, req *v1.AdminDistRuleCreateReq) (res *v1.AdminDistRuleCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRuleManage); err != nil {
		return nil, err
	}
	scopeId, err := parseID(req.ScopeId)
	if err != nil {
		return nil, err
	}
	id, err := user.NewDistributionAdminLogic().AdminRuleCreate(ctx, req.ScopeType, scopeId, req.Level1Rate, req.Level2Rate)
	if err != nil {
		return nil, err
	}
	return &v1.AdminDistRuleCreateRes{Id: fmtID(id)}, nil
}
