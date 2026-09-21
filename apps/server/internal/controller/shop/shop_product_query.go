package shop

import "ecboot/internal/model"

// productQueryFromReq C 端商品查询入参组装（列表/搜索共用; 可选 ID 空串=0）。
func productQueryFromReq(categoryId, brandId string, sort int, priceMin, priceMax string,
	page, pageSize int) (model.ProductQuery, error) {
	q := model.ProductQuery{
		Sort: sort, PriceMin: priceMin, PriceMax: priceMax,
		PageReq: model.PageReq{Page: page, PageSize: pageSize},
	}
	if categoryId != "" {
		id, err := parseID(categoryId)
		if err != nil {
			return q, err
		}
		q.CategoryId = id
	}
	if brandId != "" {
		id, err := parseID(brandId)
		if err != nil {
			return q, err
		}
		q.BrandId = id
	}
	return q, nil
}
