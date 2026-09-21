package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// CartUpdateItem 修改购物车项（数量/勾选）
func (c *ControllerV1) CartUpdateItem(ctx context.Context, req *v1.CartUpdateItemReq) (res *v1.CartUpdateItemRes, err error) {
	itemId, err := parseID(req.ItemId)
	if err != nil {
		return nil, err
	}
	checked := req.Checked
	if err = shop.NewCartLogic().UpdateItem(ctx, middleware.CtxUserIdFrom(ctx), itemId, req.Quantity, &checked); err != nil {
		return nil, err
	}
	return &v1.CartUpdateItemRes{Success: true}, nil
}
