package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminStoreDelete 删除门店(软删)
func (c *ControllerV1) AdminStoreDelete(ctx context.Context, req *v1.AdminStoreDeleteReq) (res *v1.AdminStoreDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "store:manage:delete"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.AdminDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminStoreDeleteRes{Success: true}, nil
}
