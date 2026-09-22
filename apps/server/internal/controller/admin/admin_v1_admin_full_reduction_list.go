package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// fullReductionLaddersOut/scopesOut/scopesIn 档位与范围的 api↔DTO 映射（详情/创建/更新共用, D8 同源）。
func fullReductionLaddersOut(in []model.PromotionLadder) []v1.AdminFullReductionLadder {
	out := make([]v1.AdminFullReductionLadder, 0, len(in))
	for _, l := range in {
		out = append(out, v1.AdminFullReductionLadder{Threshold: l.Threshold, Discount: l.Discount})
	}
	return out
}

func fullReductionScopesOut(in []model.PromotionScope) []v1.AdminFullReductionScope {
	out := make([]v1.AdminFullReductionScope, 0, len(in))
	for _, s := range in {
		v := ""
		if s.TargetId != 0 {
			v = fmtID(s.TargetId)
		}
		out = append(out, v1.AdminFullReductionScope{ScopeType: s.ScopeType, TargetId: v})
	}
	return out
}

func fullReductionScopesIn(in []v1.AdminFullReductionScope) ([]model.PromotionScope, error) {
	out := make([]model.PromotionScope, 0, len(in))
	for _, s := range in {
		var target int64
		if s.TargetId != "" {
			var err error
			if target, err = parseID(s.TargetId); err != nil {
				return nil, err
			}
		}
		out = append(out, model.PromotionScope{ScopeType: s.ScopeType, TargetId: target})
	}
	return out, nil
}
// AdminFullReductionList 满减活动列表。
func (c *ControllerV1) AdminFullReductionList(ctx context.Context, req *v1.AdminFullReductionListReq) (res *v1.AdminFullReductionListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFullReductionAll); err != nil {
		return nil, err
	}
	out, err := shop.NewActivityLogic().FullReductionList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminFullReductionListRes{List: make([]v1.AdminFullReductionItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminFullReductionItem{
			Id: fmtID(it.Id), Name: it.Name, StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status,
		})
	}
	return res, nil
}
