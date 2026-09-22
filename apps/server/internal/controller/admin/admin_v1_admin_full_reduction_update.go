package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminFullReductionUpdate 修改满减活动（档位/范围全量替换）。
func (c *ControllerV1) AdminFullReductionUpdate(ctx context.Context, req *v1.AdminFullReductionUpdateReq) (res *v1.AdminFullReductionUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFullReductionAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	scopes, err := fullReductionScopesIn(req.Scopes)
	if err != nil {
		return nil, err
	}
	ladders := make([]model.PromotionLadder, 0, len(req.Ladders))
	for _, l := range req.Ladders {
		ladders = append(ladders, model.PromotionLadder{Threshold: l.Threshold, Discount: l.Discount})
	}
	if err = shop.NewActivityLogic().FullReductionUpdate(ctx, id, model.PromotionActivityInput{
		Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime, Status: req.Status,
		Ladders: ladders, Scopes: scopes,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminFullReductionUpdateRes{Success: true}, nil
}
