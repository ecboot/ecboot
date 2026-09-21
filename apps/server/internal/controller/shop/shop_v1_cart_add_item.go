package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// CartAddItem 加入购物车（重复加购数量累加）
func (c *ControllerV1) CartAddItem(ctx context.Context, req *v1.CartAddItemReq) (res *v1.CartAddItemRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	skuId, err := parseID(req.SkuId)
	if err != nil {
		return nil, err
	}
	id, err := shop.NewCartLogic().AddItem(ctx, userId, skuId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &v1.CartAddItemRes{ItemId: fmtID(id)}, nil
}
