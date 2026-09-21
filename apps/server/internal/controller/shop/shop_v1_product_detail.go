package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// ProductDetail 商品详情（SKU/规格/运费概要/评价汇总; 不含成本价; 公开）
func (c *ControllerV1) ProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	// viewerUserId: 已登录会员=真实 id（用于限售/收藏等视角）, 游客=0
	v, err := shop.NewProductLogic().ProductDetail(ctx, spuId, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.ProductDetailRes{
		SpuId: fmtID(v.SpuId), SpuNo: v.SpuNo, Name: v.Name, SubTitle: v.SubTitle,
		Images: v.Images, VideoUrl: v.VideoUrl, Description: v.Description,
		SpecDefinitions: v.SpecDefinitions, Attributes: v.Attributes,
		FreightSummary: v.FreightSummary,
		Skus:           make([]v1.SkuItem, 0, len(v.Skus)),
		ReviewSummary: v1.ReviewSummary{
			Avg: v.ReviewSummary.Avg, Distribution: v.ReviewSummary.Distribution,
			Total: v.ReviewSummary.Total,
		},
	}
	for _, sk := range v.Skus {
		res.Skus = append(res.Skus, v1.SkuItem{
			SkuId: fmtID(sk.SkuId), SkuNo: sk.SkuNo, Specs: sk.Specs,
			Price: sk.Price, LinePrice: sk.LinePrice, Sellable: sk.Sellable,
		})
	}
	return res, nil
}
