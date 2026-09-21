package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// MessageReadAll 全部已读
func (c *ControllerV1) MessageReadAll(ctx context.Context, req *v1.MessageReadAllReq) (res *v1.MessageReadAllRes, err error) {
	if err = user.MarkAllRead(ctx, middleware.CtxUserIdFrom(ctx)); err != nil {
		return nil, err
	}
	return &v1.MessageReadAllRes{Success: true}, nil
}
