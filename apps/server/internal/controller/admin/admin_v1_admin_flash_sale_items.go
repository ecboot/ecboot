package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminFlashSaleItems 秒杀场次商品全量替换（秒杀价/限量/限购; 已售行不可移除）。
func (c *ControllerV1) AdminFlashSaleItems(ctx context.Context, req *v1.AdminFlashSaleItemsReq) (res *v1.AdminFlashSaleItemsRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFlashSaleAll); err != nil {
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
	if err = shop.NewActivityLogic().FlashSaleSetItems(ctx, id, items); err != nil {
		return nil, err
	}
	return &v1.AdminFlashSaleItemsRes{Success: true}, nil
}
