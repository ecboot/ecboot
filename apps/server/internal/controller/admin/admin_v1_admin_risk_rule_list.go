package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/system"
)

// AdminRiskRuleList 风控规则列表。
func (c *ControllerV1) AdminRiskRuleList(ctx context.Context, req *v1.AdminRiskRuleListReq) (res *v1.AdminRiskRuleListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermRiskRuleRead); err != nil {
		return nil, err
	}
	out, err := system.NewRiskAdminLogic().AdminRuleList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminRiskRuleListRes{List: make([]v1.AdminRiskRuleItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminRiskRuleItem{
			Id: fmtID(it.Id), Name: it.Name, RuleType: it.RuleType,
			ConditionExpr: it.ConditionExpr, Action: it.Action, Status: it.Status,
		})
	}
	return res, nil
}
