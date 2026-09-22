package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminGroupBuyCreate 创建拼团活动。
func (c *ControllerV1) AdminGroupBuyCreate(ctx context.Context, req *v1.AdminGroupBuyCreateReq) (res *v1.AdminGroupBuyCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionGroupBuyAll); err != nil {
		return nil, err
	}
	spuid, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	id, err := shop.NewActivityLogic().GroupBuyCreate(ctx, model.GroupBuyInput{
		Name: req.Name, SpuId: spuid, GroupSize: req.GroupSize, PerLimit: req.PerLimit,
		StartTime: req.StartTime, EndTime: req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminGroupBuyCreateRes{Id: fmtID(id)}, nil
}
