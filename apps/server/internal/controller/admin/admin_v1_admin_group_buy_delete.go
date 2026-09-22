package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)


// AdminGroupBuyDelete 软删拼团活动（C 端列表即时消失）。
func (c *ControllerV1) AdminGroupBuyDelete(ctx context.Context, req *v1.AdminGroupBuyDeleteReq) (res *v1.AdminGroupBuyDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionGroupBuyAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().GroupBuyDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminGroupBuyDeleteRes{Success: true}, nil
}
