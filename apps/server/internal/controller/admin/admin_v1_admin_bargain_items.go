package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminBargainItems 砍价场次商品全量替换（起始价/底价/最大刀数）。
func (c *ControllerV1) AdminBargainItems(ctx context.Context, req *v1.AdminBargainItemsReq) (res *v1.AdminBargainItemsRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionBargainAll); err != nil {
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
	if err = shop.NewActivityLogic().BargainSetItems(ctx, id, items); err != nil {
		return nil, err
	}
	return &v1.AdminBargainItemsRes{Success: true}, nil
}
