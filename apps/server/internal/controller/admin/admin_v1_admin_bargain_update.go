package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminBargainUpdate 修改砍价活动（status 显式启停）。
func (c *ControllerV1) AdminBargainUpdate(ctx context.Context, req *v1.AdminBargainUpdateReq) (res *v1.AdminBargainUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionBargainAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().BargainUpdate(ctx, id, model.BargainActivityInput{
		Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminBargainUpdateRes{Success: true}, nil
}
