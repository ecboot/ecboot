package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminSkuCreate 新增SKU
func (c *ControllerV1) AdminSkuCreate(ctx context.Context, req *v1.AdminSkuCreateReq) (res *v1.AdminSkuCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:sku:create"); err != nil {
		return nil, err
	}
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	id, skuNo, err := shop.AdminSkuCreateWithNo(ctx, spuId, model.SkuInput{
		Specs: req.Specs, Price: req.Price, LinePrice: req.LinePrice, CostPrice: req.CostPrice,
		Image: req.Image, Weight: req.Weight, Barcode: req.Barcode,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminSkuCreateRes{SkuId: fmtID(id), SkuNo: skuNo}, nil
}
