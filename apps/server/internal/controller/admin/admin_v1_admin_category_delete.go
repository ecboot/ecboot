package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminCategoryDelete 删除分类(软删,有商品禁删)
func (c *ControllerV1) AdminCategoryDelete(ctx context.Context, req *v1.AdminCategoryDeleteReq) (res *v1.AdminCategoryDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:category:delete"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminCategoryDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminCategoryDeleteRes{Success: true}, nil
}
