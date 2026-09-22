package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminAssistUpdate 修改助力活动（status 显式启停; 奖励载体不可改——api 契约无字段）。
func (c *ControllerV1) AdminAssistUpdate(ctx context.Context, req *v1.AdminAssistUpdateReq) (res *v1.AdminAssistUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionAssistAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().AssistUpdate(ctx, id, model.AssistActivityInput{
		Name: req.Name, RequiredCount: req.RequiredCount, PerLimit: req.PerLimit,
		StartTime: req.StartTime, EndTime: req.EndTime, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminAssistUpdateRes{Success: true}, nil
}
