package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// ProductList 商品列表（在售; 筛选与排序; 公开）
func (c *ControllerV1) ProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	q, err := productQueryFromReq(req.CategoryId, req.BrandId, req.Sort, req.PriceMin, req.PriceMax,
		req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewProductLogic().Products(ctx, q)
	if err != nil {
		return nil, err
	}
	res = &v1.ProductListRes{List: make([]v1.ProductItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.ProductItem{
			SpuId: fmtID(it.SpuId), SpuName: it.Name, Image: it.Image,
			PriceRange: it.PriceRange, SaleCount: it.SaleCount,
		})
	}
	return res, nil
}
