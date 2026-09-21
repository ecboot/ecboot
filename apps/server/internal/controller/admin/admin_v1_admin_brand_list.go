package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminBrandList 品牌列表
func (c *ControllerV1) AdminBrandList(ctx context.Context, req *v1.AdminBrandListReq) (res *v1.AdminBrandListRes, err error) {
	out, err := shop.NewProductLogic().AdminBrandList(ctx, req.Status,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminBrandListRes{List: make([]v1.AdminBrandItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminBrandItem{
			Id: fmtID(it.Id), Name: it.Name, Logo: it.Logo,
			Description: it.Description, Sort: it.Sort, Status: it.Status,
		})
	}
	return res, nil
}
