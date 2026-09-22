package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// GroupBuyList 拼团活动列表（公开; 只出进行中, 按结束时间升序）
func (c *ControllerV1) GroupBuyList(ctx context.Context, req *v1.GroupBuyListReq) (res *v1.GroupBuyListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicGroupBuys(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.GroupBuyListRes{List: groupBuyItems(out.List)}
	res.Total = out.Total
	return res, nil
}
