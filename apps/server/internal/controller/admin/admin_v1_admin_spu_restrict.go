package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminSpuRestrict 限售区域设置
func (c *ControllerV1) AdminSpuRestrict(ctx context.Context, req *v1.AdminSpuRestrictReq) (res *v1.AdminSpuRestrictRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:spu:update"); err != nil {
		return nil, err
	}
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminProductRestrict(ctx, spuId, req.SaleRestrictCodes); err != nil {
		return nil, err
	}
	return &v1.AdminSpuRestrictRes{Success: true}, nil
}
