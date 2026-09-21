package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// ProfileUpdate 修改个人资料（未传字段保持不变）
func (c *ControllerV1) ProfileUpdate(ctx context.Context, req *v1.ProfileUpdateReq) (res *v1.ProfileUpdateRes, err error) {
	if err = user.ProfileUpdate(ctx, middleware.CtxUserIdFrom(ctx), model.ProfileUpdateInput{
		Nickname: req.Nickname, Avatar: req.Avatar, Gender: req.Gender,
	}); err != nil {
		return nil, err
	}
	return &v1.ProfileUpdateRes{Success: true}, nil
}
