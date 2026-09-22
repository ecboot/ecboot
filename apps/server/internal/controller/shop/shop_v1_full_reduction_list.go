package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// FullReductionList 满减活动列表（公开; 只出进行中）
func (c *ControllerV1) FullReductionList(ctx context.Context, req *v1.FullReductionListReq) (res *v1.FullReductionListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicFullReductions(ctx, model.PageReq{Page: 1, PageSize: 100})
	if err != nil {
		return nil, err
	}
	res = &v1.FullReductionListRes{List: fullReductionItems(out.List)}
	return res, nil
}
