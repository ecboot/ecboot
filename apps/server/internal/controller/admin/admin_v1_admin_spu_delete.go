package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminSpuDelete 删除商品(软删)
func (c *ControllerV1) AdminSpuDelete(ctx context.Context, req *v1.AdminSpuDeleteReq) (res *v1.AdminSpuDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:spu:delete"); err != nil {
		return nil, err
	}
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminProductDelete(ctx, spuId); err != nil {
		return nil, err
	}
	return &v1.AdminSpuDeleteRes{Success: true}, nil
}
