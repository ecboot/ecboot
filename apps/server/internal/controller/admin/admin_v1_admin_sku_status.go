package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminSkuStatus SKU启停
func (c *ControllerV1) AdminSkuStatus(ctx context.Context, req *v1.AdminSkuStatusReq) (res *v1.AdminSkuStatusRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:sku:update"); err != nil {
		return nil, err
	}
	skuId, err := parseID(req.SkuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminSkuStatus(ctx, skuId, req.Status); err != nil {
		return nil, err
	}
	return &v1.AdminSkuStatusRes{Success: true}, nil
}
