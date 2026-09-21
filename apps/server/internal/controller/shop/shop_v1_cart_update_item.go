package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// CartUpdateItem 修改购物车项（数量/勾选）
func (c *ControllerV1) CartUpdateItem(ctx context.Context, req *v1.CartUpdateItemReq) (res *v1.CartUpdateItemRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	itemId, err := parseID(req.ItemId)
	if err != nil {
		return nil, err
	}
	// 评审 I9: 直传 api 的三态指针——原 `checked := req.Checked; &checked` 恒非 nil,
	// 使"只改数量"的请求静默把勾选清掉（结算金额随之变化）。
	if err = shop.NewCartLogic().UpdateItem(ctx, userId, itemId, req.Quantity, req.Checked); err != nil {
		return nil, err
	}
	return &v1.CartUpdateItemRes{Success: true}, nil
}
