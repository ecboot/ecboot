package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminFlashSaleUpdate 修改秒杀活动（status 显式启停）。
func (c *ControllerV1) AdminFlashSaleUpdate(ctx context.Context, req *v1.AdminFlashSaleUpdateReq) (res *v1.AdminFlashSaleUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFlashSaleAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().FlashSaleUpdate(ctx, id, model.ActivityTimeInput{
		Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminFlashSaleUpdateRes{Success: true}, nil
}
