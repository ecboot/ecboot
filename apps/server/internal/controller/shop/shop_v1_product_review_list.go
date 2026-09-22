package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// ProductReviewList 商品评价列表（只出审核通过; 含汇总与星级筛选）
func (c *ControllerV1) ProductReviewList(ctx context.Context, req *v1.ProductReviewListReq) (res *v1.ProductReviewListRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	summary, page, err := shop.NewReviewLogic().ProductList(ctx, spuId, req.Score,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.ProductReviewListRes{
		Summary: v1.ReviewSummary{
			Avg: summary.Avg, Distribution: summary.Distribution, Total: summary.Total,
		},
		List: make([]v1.ReviewItem, 0, len(page.List)),
	}
	res.Total = page.Total
	for _, it := range page.List {
		res.List = append(res.List, v1.ReviewItem{
			ReviewId:  fmtID(it.ReviewId),
			User:      it.User,
			Score:     it.Score,
			Content:   it.Content,
			Images:    it.Images,
			Specs:     it.Specs,
			Reply:     it.Reply,
			Extra:     it.Extra,
			CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
