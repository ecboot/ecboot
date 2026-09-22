package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)


// AdminAssistDelete 软删助力活动（C 端列表即时消失）。
func (c *ControllerV1) AdminAssistDelete(ctx context.Context, req *v1.AdminAssistDeleteReq) (res *v1.AdminAssistDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionAssistAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewActivityLogic().AssistDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminAssistDeleteRes{Success: true}, nil
}
