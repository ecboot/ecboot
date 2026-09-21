package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminSkuDelete 删除SKU(软删)
func (c *ControllerV1) AdminSkuDelete(ctx context.Context, req *v1.AdminSkuDeleteReq) (res *v1.AdminSkuDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:sku:delete"); err != nil {
		return nil, err
	}
	skuId, err := parseID(req.SkuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminSkuDelete(ctx, skuId); err != nil {
		return nil, err
	}
	return &v1.AdminSkuDeleteRes{Success: true}, nil
}
