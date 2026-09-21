package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// ProfileDetail 个人资料（手机号脱敏）
func (c *ControllerV1) ProfileDetail(ctx context.Context, req *v1.ProfileDetailReq) (res *v1.ProfileDetailRes, err error) {
	d, err := user.ProfileDetail(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.ProfileDetailRes{
		UserId: d.UserId, Nickname: d.Nickname, Avatar: d.Avatar, Gender: d.Gender,
		Phone: d.PhoneMasked, HasPassword: d.HasPassword,
		GrowthValue: d.GrowthValue, Level: int(d.Level), LevelName: d.LevelName,
	}, nil
}
