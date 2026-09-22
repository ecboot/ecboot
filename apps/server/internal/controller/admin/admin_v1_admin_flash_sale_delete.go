package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)


// AdminFlashSaleDelete 软删秒杀活动（C 端列表即时消失）。
func (c *ControllerV1) AdminFlashSaleDelete(ctx context.Context, req *v1.AdminFlashSaleDeleteReq) (res *v1.AdminFlashSaleDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFlashSaleAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().FlashSaleDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminFlashSaleDeleteRes{Success: true}, nil
}
