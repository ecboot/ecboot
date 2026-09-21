package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// NotifyPreferenceSet 设置通知偏好（幂等）
func (c *ControllerV1) NotifyPreferenceSet(ctx context.Context, req *v1.NotifyPreferenceSetReq) (res *v1.NotifyPreferenceSetRes, err error) {
	prefs := make([]model.NotifyPreference, 0, len(req.List))
	for _, p := range req.List {
		prefs = append(prefs, model.NotifyPreference{Channel: p.Channel, Enabled: p.Enabled})
	}
	if err = user.SetPreferences(ctx, middleware.CtxUserIdFrom(ctx), prefs); err != nil {
		return nil, err
	}
	return &v1.NotifyPreferenceSetRes{Success: true}, nil
}
