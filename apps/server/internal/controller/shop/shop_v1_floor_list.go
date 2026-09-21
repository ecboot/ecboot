package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// FloorList 楼层内容（商品楼层含装配摘要; 公开）
func (c *ControllerV1) FloorList(ctx context.Context, req *v1.FloorListReq) (res *v1.FloorListRes, err error) {
	list, err := shop.PublicFloors(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.FloorListRes{List: make([]v1.FloorItem, 0, len(list))}
	for _, it := range list {
		item := v1.FloorItem{
			FloorId:   fmtID(it.FloorId),
			FloorType: it.FloorType,
			Title:     it.Title,
			Config:    it.Config,
			Products:  make([]v1.FloorProduct, 0, len(it.Products)),
		}
		for _, p := range it.Products {
			item.Products = append(item.Products, v1.FloorProduct{
				SpuId: fmtID(p.SpuId), Name: p.Name, Image: p.Image, Price: p.Price,
			})
		}
		res.List = append(res.List, item)
	}
	return res, nil
}
