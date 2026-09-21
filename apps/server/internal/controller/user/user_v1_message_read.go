package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// MessageRead 标记已读（归属校验）
func (c *ControllerV1) MessageRead(ctx context.Context, req *v1.MessageReadReq) (res *v1.MessageReadRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.MarkRead(ctx, middleware.CtxUserIdFrom(ctx), id); err != nil {
		return nil, err
	}
	return &v1.MessageReadRes{Success: true}, nil
}
