package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminAssistCreate 创建助力活动（D5: 券奖励校验存在; 积分奖励存 config）。
func (c *ControllerV1) AdminAssistCreate(ctx context.Context, req *v1.AdminAssistCreateReq) (res *v1.AdminAssistCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionAssistAll); err != nil {
		return nil, err
	}
	var rewardRef int64
	if req.RewardRef != "" {
		if rewardRef, err = parseID(req.RewardRef); err != nil {
			return nil, err
		}
	}
	id, err := shop.NewActivityLogic().AssistCreate(ctx, model.AssistActivityInput{
		Name: req.Name, RewardType: req.RewardType, RewardRef: rewardRef, PointAmount: req.PointAmount,
		RequiredCount: req.RequiredCount, PerLimit: req.PerLimit,
		StartTime: req.StartTime, EndTime: req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminAssistCreateRes{Id: fmtID(id)}, nil
}
