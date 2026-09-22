package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminFullReductionCreate 创建满减活动（档位+范围嵌套; 空范围=全场）。
func (c *ControllerV1) AdminFullReductionCreate(ctx context.Context, req *v1.AdminFullReductionCreateReq) (res *v1.AdminFullReductionCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFullReductionAll); err != nil {
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
	id, err := shop.NewActivityLogic().FullReductionCreate(ctx, model.PromotionActivityInput{
		Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime,
		Ladders: ladders, Scopes: scopes,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminFullReductionCreateRes{Id: fmtID(id)}, nil
}
