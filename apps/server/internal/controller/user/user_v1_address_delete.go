package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AddressDelete 删除收货地址（归属校验）
func (c *ControllerV1) AddressDelete(ctx context.Context, req *v1.AddressDeleteReq) (res *v1.AddressDeleteRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.AddressDelete(ctx, middleware.CtxUserIdFrom(ctx), id); err != nil {
		return nil, err
	}
	return &v1.AddressDeleteRes{Success: true}, nil
}
