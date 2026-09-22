package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminGroupBuyItems 拼团场次商品全量替换（成团价必填）。
func (c *ControllerV1) AdminGroupBuyItems(ctx context.Context, req *v1.AdminGroupBuyItemsReq) (res *v1.AdminGroupBuyItemsRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionGroupBuyAll); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	items := make([]model.ActivitySkuInput, 0, len(req.Items))
	for _, it := range req.Items {
		var skuId int64
		if skuId, err = parseID(it.SkuId); err != nil {
			return nil, err
		}
		items = append(items, model.ActivitySkuInput{SkuId: skuId, GroupPrice: it.GroupPrice,
			FlashPrice: it.FlashPrice, StockCount: it.StockCount, PerLimit: it.PerLimit,
			OriginalPrice: it.OriginalPrice, FloorPrice: it.FloorPrice, MaxCutCount: it.MaxCutCount})
	}
	if err = shop.NewActivityLogic().GroupBuySetItems(ctx, id, items); err != nil {
		return nil, err
	}
	return &v1.AdminGroupBuyItemsRes{Success: true}, nil
}
