package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// CartDetail 购物车详情
func (c *ControllerV1) CartDetail(ctx context.Context, req *v1.CartDetailReq) (res *v1.CartDetailRes, err error) {
	view, err := shop.NewCartLogic().Detail(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.CartDetailRes{List: make([]v1.CartItem, 0, len(view.Items))}
	for _, it := range view.Items {
		res.List = append(res.List, v1.CartItem{
			ItemId: fmtID(it.ItemId), SkuId: fmtID(it.SkuId), SpuName: it.SpuName,
			Specs: it.Specs, Image: it.Image, Price: it.Price,
			Sellable: it.Sellable, Quantity: it.Quantity, Checked: it.Checked,
		})
	}
	return res, nil
}
