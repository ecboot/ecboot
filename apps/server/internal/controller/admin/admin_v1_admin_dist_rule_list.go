package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// AdminDistRuleList 佣金规则列表。
func (c *ControllerV1) AdminDistRuleList(ctx context.Context, req *v1.AdminDistRuleListReq) (res *v1.AdminDistRuleListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermDistributionRuleRead); err != nil {
		return nil, err
	}
	out, err := user.NewDistributionAdminLogic().AdminRuleList(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminDistRuleListRes{List: make([]v1.AdminDistRuleItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminDistRuleItem{
			Id: fmtID(it.Id), ScopeType: it.ScopeType, ScopeDesc: it.ScopeDesc,
			Level1Rate: it.Level1Rate, Level2Rate: it.Level2Rate, Status: it.Status,
		})
	}
	return res, nil
}
