package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// FlashSaleList 秒杀活动列表（公开; 进行中 + 预告）
func (c *ControllerV1) FlashSaleList(ctx context.Context, req *v1.FlashSaleListReq) (res *v1.FlashSaleListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicFlashSales(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.FlashSaleListRes{List: flashSaleItems(out.List)}
	res.Total = out.Total
	return res, nil
}
