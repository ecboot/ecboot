package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminFullReductionDetail 满减详情（档位升序 + 范围回读）。
func (c *ControllerV1) AdminFullReductionDetail(ctx context.Context, req *v1.AdminFullReductionDetailReq) (res *v1.AdminFullReductionDetailRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFullReductionAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	d, err := shop.NewActivityLogic().FullReductionDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminFullReductionDetailRes{
		Id: fmtID(d.Id), Name: d.Name, StartTime: d.StartTime, EndTime: d.EndTime, Status: d.Status,
		Ladders: fullReductionLaddersOut(d.Ladders), Scopes: fullReductionScopesOut(d.Scopes),
	}, nil
}
