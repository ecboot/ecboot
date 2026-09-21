package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminSkuUpdate 修改SKU
func (c *ControllerV1) AdminSkuUpdate(ctx context.Context, req *v1.AdminSkuUpdateReq) (res *v1.AdminSkuUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:sku:update"); err != nil {
		return nil, err
	}
	skuId, err := parseID(req.SkuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminSkuUpdate(ctx, skuId, model.SkuInput{
		Specs: req.Specs, Price: req.Price, LinePrice: req.LinePrice, CostPrice: req.CostPrice,
		Image: req.Image, Weight: req.Weight, Barcode: req.Barcode,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminSkuUpdateRes{Success: true}, nil
}
