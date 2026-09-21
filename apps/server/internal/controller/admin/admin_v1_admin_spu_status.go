package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminSpuStatus 商品上下架（无启用 SKU 禁上架）
func (c *ControllerV1) AdminSpuStatus(ctx context.Context, req *v1.AdminSpuStatusReq) (res *v1.AdminSpuStatusRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:spu:update"); err != nil {
		return nil, err
	}
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminProductStatus(ctx, spuId, req.Status); err != nil {
		return nil, err
	}
	return &v1.AdminSpuStatusRes{Success: true}, nil
}
