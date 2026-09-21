package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AddressSetDefault 设默认地址（归属校验; 同事务清其他, 默认唯一）
func (c *ControllerV1) AddressSetDefault(ctx context.Context, req *v1.AddressSetDefaultReq) (res *v1.AddressSetDefaultRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.AddressSetDefault(ctx, middleware.CtxUserIdFrom(ctx), id); err != nil {
		return nil, err
	}
	return &v1.AddressSetDefaultRes{Success: true}, nil
}
