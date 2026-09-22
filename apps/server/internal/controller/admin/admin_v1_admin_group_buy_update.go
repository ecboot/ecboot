package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminGroupBuyUpdate 修改拼团活动（status 显式启停）。
func (c *ControllerV1) AdminGroupBuyUpdate(ctx context.Context, req *v1.AdminGroupBuyUpdateReq) (res *v1.AdminGroupBuyUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionGroupBuyAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().GroupBuyUpdate(ctx, id, model.GroupBuyInput{
		Name: req.Name, GroupSize: req.GroupSize, PerLimit: req.PerLimit,
		StartTime: req.StartTime, EndTime: req.EndTime, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminGroupBuyUpdateRes{Success: true}, nil
}
