package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminFullReductionDelete 软删满减活动（C 端列表即时消失）。
func (c *ControllerV1) AdminFullReductionDelete(ctx context.Context, req *v1.AdminFullReductionDeleteReq) (res *v1.AdminFullReductionDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFullReductionAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().FullReductionDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminFullReductionDeleteRes{Success: true}, nil
}
