package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// BrandList 品牌列表（启用; 公开）
func (c *ControllerV1) BrandList(ctx context.Context, req *v1.BrandListReq) (res *v1.BrandListRes, err error) {
	out, err := shop.NewProductLogic().Brands(ctx, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.BrandListRes{List: make([]v1.BrandItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.BrandItem{Id: fmtID(it.Id), Name: it.Name, Logo: it.Logo})
	}
	return res, nil
}
