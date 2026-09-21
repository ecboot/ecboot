package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// CartRemoveItem 移除购物车项
func (c *ControllerV1) CartRemoveItem(ctx context.Context, req *v1.CartRemoveItemReq) (res *v1.CartRemoveItemRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	itemId, err := parseID(req.ItemId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewCartLogic().RemoveItem(ctx, userId, itemId); err != nil {
		return nil, err
	}
	return &v1.CartRemoveItemRes{Success: true}, nil
}
