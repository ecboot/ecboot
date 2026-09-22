package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)


// AdminBargainDelete 软删砍价活动（C 端列表即时消失）。
func (c *ControllerV1) AdminBargainDelete(ctx context.Context, req *v1.AdminBargainDeleteReq) (res *v1.AdminBargainDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionBargainAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().BargainDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminBargainDeleteRes{Success: true}, nil
}
