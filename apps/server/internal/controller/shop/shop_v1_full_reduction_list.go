package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// FullReductionList 满减活动列表（公开; 只出进行中; 可按商品过滤范围命中 + 分页）
// I3（015 修复轮）: 原控制器硬编码 `PageSize: 100` 且丢弃 spuId——FR-018 的分页
// 与 FR-016/018 的范围命中均不成立, 现按契约透传。
func (c *ControllerV1) FullReductionList(ctx context.Context, req *v1.FullReductionListReq) (res *v1.FullReductionListRes, err error) {
	out, err := shop.NewMarketingLogic().PublicFullReductions(ctx, spuIdOf(req.SpuId), model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.FullReductionListRes{List: fullReductionItems(out.List)}
	res.Total = out.Total
	return res, nil
}
