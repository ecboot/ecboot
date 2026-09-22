package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminFlashSaleCreate 创建秒杀活动。
func (c *ControllerV1) AdminFlashSaleCreate(ctx context.Context, req *v1.AdminFlashSaleCreateReq) (res *v1.AdminFlashSaleCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFlashSaleAll); err != nil {
		return nil, err
	}
	id, err := shop.NewActivityLogic().FlashSaleCreate(ctx, model.ActivityTimeInput{
		Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminFlashSaleCreateRes{Id: fmtID(id)}, nil
}
