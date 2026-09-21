package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// ProductSearch 关键词搜索（公开）
func (c *ControllerV1) ProductSearch(ctx context.Context, req *v1.ProductSearchReq) (res *v1.ProductSearchRes, err error) {
	q, err := productQueryFromReq(req.CategoryId, req.BrandId, req.Sort, "", "",
		req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewProductLogic().Search(ctx, req.Keyword, q)
	if err != nil {
		return nil, err
	}
	res = &v1.ProductSearchRes{List: make([]v1.ProductItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.ProductItem{
			SpuId: fmtID(it.SpuId), SpuName: it.Name, Image: it.Image,
			PriceRange: it.PriceRange, SaleCount: it.SaleCount,
		})
	}
	return res, nil
}
