package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminBargainCreate 创建砍价活动。
func (c *ControllerV1) AdminBargainCreate(ctx context.Context, req *v1.AdminBargainCreateReq) (res *v1.AdminBargainCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionBargainAll); err != nil {
		return nil, err
	}
	spuid, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	id, err := shop.NewActivityLogic().BargainCreate(ctx, model.BargainActivityInput{
		Name: req.Name, SpuId: spuid, StartTime: req.StartTime, EndTime: req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminBargainCreateRes{Id: fmtID(id)}, nil
}
