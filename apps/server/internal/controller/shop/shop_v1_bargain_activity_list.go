package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// BargainActivityList 砍价活动列表（公开; 只出进行中）
func (c *ControllerV1) BargainActivityList(ctx context.Context, req *v1.BargainActivityListReq) (res *v1.BargainActivityListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicBargains(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.BargainActivityListRes{List: bargainItems(out.List)}
	res.Total = out.Total
	return res, nil
}
