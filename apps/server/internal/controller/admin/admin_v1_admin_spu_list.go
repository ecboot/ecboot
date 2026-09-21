package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminSpuList 商品SPU列表（全状态）
func (c *ControllerV1) AdminSpuList(ctx context.Context, req *v1.AdminSpuListReq) (res *v1.AdminSpuListRes, err error) {
	var catId int64
	if req.CategoryId != "" {
		if catId, err = parseID(req.CategoryId); err != nil {
			return nil, err
		}
	}
	out, err := shop.NewProductLogic().AdminProductList(ctx, model.AdminProductQuery{
		Status: req.Status, CategoryId: catId, Keyword: req.Keyword,
		PageReq: model.PageReq{Page: req.Page, PageSize: req.PageSize},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminSpuListRes{List: make([]v1.AdminSpuItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminSpuItem{
			SpuId: fmtID(it.SpuId), SpuNo: it.SpuNo, Name: it.Name,
			CategoryId: fmtID(it.CategoryId), BrandId: fmtID(it.BrandId),
			Status: it.Status, SaleCount: it.SaleCount, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
