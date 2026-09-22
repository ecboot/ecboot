package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AssistList 助力活动列表（公开; 只出进行中）
func (c *ControllerV1) AssistList(ctx context.Context, req *v1.AssistListReq) (res *v1.AssistListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicAssists(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AssistListRes{List: assistItems(out.List)}
	res.Total = out.Total
	return res, nil
}
