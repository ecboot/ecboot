package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// NotifyPreferenceGet 通知偏好（未设置=默认全开）
func (c *ControllerV1) NotifyPreferenceGet(ctx context.Context, req *v1.NotifyPreferenceGetReq) (res *v1.NotifyPreferenceGetRes, err error) {
	prefs, err := user.Preferences(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.NotifyPreferenceGetRes{List: make([]v1.NotifyPreference, 0, len(prefs))}
	for _, p := range prefs {
		res.List = append(res.List, v1.NotifyPreference{Channel: p.Channel, Enabled: p.Enabled})
	}
	return res, nil
}
