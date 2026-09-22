package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// Index 首页聚合（公开）: 轮播 + 楼层 + 五类活动入口(各至多 3 条) + 可领券; 任一块为空返回空数组。
func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	agg, err := shop.NewMarketingLogic().Index(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.IndexRes{
		Banners:        make([]v1.BannerItem, 0, len(agg.Banners)),
		Floors:         make([]v1.FloorItem, 0, len(agg.Floors)),
		FlashSales:     make([]v1.FlashSaleItem, 0, len(agg.FlashSales)),
		GroupBuys:      make([]v1.GroupBuyItem, 0, len(agg.GroupBuys)),
		Bargains:       make([]v1.BargainActivityItem, 0, len(agg.Bargains)),
		Assists:        make([]v1.AssistActivityItem, 0, len(agg.Assists)),
		FullReductions: make([]v1.FullReductionItem, 0, len(agg.FullReductions)),
		Coupons:        make([]v1.IndexCoupon, 0, len(agg.Coupons)),
	}
	for _, b := range agg.Banners {
		res.Banners = append(res.Banners, v1.BannerItem{Id: fmtID(b.Id), ImageUrl: b.ImageUrl, LinkUrl: b.LinkUrl})
	}
	for _, fl := range agg.Floors {
		item := v1.FloorItem{
			FloorId: fmtID(fl.FloorId), FloorType: fl.FloorType, Title: fl.Title, Config: fl.Config,
			Products: make([]v1.FloorProduct, 0, len(fl.Products)),
		}
		for _, p := range fl.Products {
			item.Products = append(item.Products, v1.FloorProduct{
				SpuId: fmtID(p.SpuId), Name: p.Name, Image: p.Image, Price: p.Price,
			})
		}
		res.Floors = append(res.Floors, item)
	}
	res.FlashSales = flashSaleItems(agg.FlashSales)
	res.GroupBuys = groupBuyItems(agg.GroupBuys)
	res.Bargains = bargainItems(agg.Bargains)
	res.Assists = assistItems(agg.Assists)
	res.FullReductions = fullReductionItems(agg.FullReductions)
	for _, cp := range agg.Coupons {
		res.Coupons = append(res.Coupons, v1.IndexCoupon{
			Id: fmtID(cp.Id), Name: cp.Name, Threshold: cp.Threshold, Discount: cp.Discount, ValidDesc: cp.ValidDesc,
		})
	}
	return res, nil
}

// ---- 分块映射（列表端点与首页入口共用同一映射, 避免两处漂移）----

func flashSaleItems(list []model.PublicFlashSaleItem) []v1.FlashSaleItem {
	out := make([]v1.FlashSaleItem, 0, len(list))
	for _, it := range list {
		row := v1.FlashSaleItem{ActivityId: fmtID(it.ActivityId), Name: it.Name,
			StartTime: it.StartTime, EndTime: it.EndTime, Items: make([]v1.FlashSaleSku, 0, len(it.Items))}
		for _, b := range it.Items {
			row.Items = append(row.Items, v1.FlashSaleSku{
				SkuId: fmtID(b.SkuId), FlashPrice: b.Price, StockRemain: b.StockRemain, PerLimit: b.PerLimit,
			})
		}
		out = append(out, row)
	}
	return out
}

func groupBuyItems(list []model.PublicGroupBuyItem) []v1.GroupBuyItem {
	out := make([]v1.GroupBuyItem, 0, len(list))
	for _, it := range list {
		row := v1.GroupBuyItem{ActivityId: fmtID(it.ActivityId), Name: it.Name, SpuId: fmtID(it.SpuId),
			SpuName: it.SpuName, Image: it.Image, GroupSize: it.GroupSize, EndTime: it.EndTime,
			Items: make([]v1.GroupBuySku, 0, len(it.Items))}
		for _, b := range it.Items {
			row.Items = append(row.Items, v1.GroupBuySku{SkuId: fmtID(b.SkuId), GroupPrice: b.Price})
		}
		out = append(out, row)
	}
	return out
}

func bargainItems(list []model.PublicBargainItem) []v1.BargainActivityItem {
	out := make([]v1.BargainActivityItem, 0, len(list))
	for _, it := range list {
		row := v1.BargainActivityItem{ActivityId: fmtID(it.ActivityId), Name: it.Name, SpuId: fmtID(it.SpuId),
			SpuName: it.SpuName, Image: it.Image, EndTime: it.EndTime,
			Items: make([]v1.BargainSku, 0, len(it.Items))}
		for _, b := range it.Items {
			row.Items = append(row.Items, v1.BargainSku{
				ItemId: fmtID(b.ItemId), SkuId: fmtID(b.SkuId),
				OriginalPrice: b.OriginalPrice, FloorPrice: b.FloorPrice, MaxCutCount: b.MaxCutCount,
			})
		}
		out = append(out, row)
	}
	return out
}

func assistItems(list []model.PublicAssistItem) []v1.AssistActivityItem {
	out := make([]v1.AssistActivityItem, 0, len(list))
	for _, it := range list {
		desc := "优惠券"
		if it.RewardType == 2 {
			desc = "积分"
		}
		out = append(out, v1.AssistActivityItem{
			ActivityId: fmtID(it.ActivityId), Name: it.Name, RewardDesc: desc,
			RequiredCount: it.RequiredCount, EndTime: it.EndTime,
		})
	}
	return out
}

func fullReductionItems(list []model.PublicFullReductionItem) []v1.FullReductionItem {
	out := make([]v1.FullReductionItem, 0, len(list))
	for _, fr := range list {
		row := v1.FullReductionItem{ActivityId: fmtID(fr.ActivityId), Name: fr.Name,
			Ladders: make([]v1.FullReductionLadder, 0, len(fr.Ladders))}
		for _, l := range fr.Ladders {
			row.Ladders = append(row.Ladders, v1.FullReductionLadder{Threshold: l.Threshold, Discount: l.Discount})
		}
		out = append(out, row)
	}
	return out
}
