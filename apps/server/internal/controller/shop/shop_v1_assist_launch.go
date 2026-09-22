package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AssistLaunch 发起助力（会员; 受每人可发起次数限制）
func (c *ControllerV1) AssistLaunch(ctx context.Context, req *v1.AssistLaunchReq) (res *v1.AssistLaunchRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	activityId, err := parseID(req.ActivityId)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewAssistLogic().Launch(ctx, userId, activityId)
	if err != nil {
		return nil, err
	}
	return &v1.AssistLaunchRes{RecordId: fmtID(out.RecordId)}, nil
}
