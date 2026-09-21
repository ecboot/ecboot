package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminBrandDelete 删除品牌(软删)
func (c *ControllerV1) AdminBrandDelete(ctx context.Context, req *v1.AdminBrandDeleteReq) (res *v1.AdminBrandDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:brand:delete"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminBrandDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminBrandDeleteRes{Success: true}, nil
}
