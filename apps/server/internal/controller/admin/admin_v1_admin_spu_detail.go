package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/shop"
)

// AdminSpuDetail 商品详情（管理视角,含成本价）
func (c *ControllerV1) AdminSpuDetail(ctx context.Context, req *v1.AdminSpuDetailReq) (res *v1.AdminSpuDetailRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	v, err := shop.NewProductLogic().AdminProductDetail(ctx, spuId)
	if err != nil {
		return nil, err
	}
	res = &v1.AdminSpuDetailRes{
		SpuId: fmtID(v.SpuId), SpuNo: v.SpuNo, Name: v.Name, SubTitle: v.SubTitle,
		CategoryId: fmtID(v.CategoryId), BrandId: fmtID(v.BrandId),
		FreightTemplateId: fmtID(v.FreightTemplateId), Images: v.Images, VideoUrl: v.VideoUrl,
		Description: v.Description, SpecDefinitions: v.SpecDefinitions, Attributes: v.Attributes,
		SaleRestrictCodes: v.SaleRestrictCodes, Status: v.Status,
		Skus: make([]v1.AdminSkuAdminItem, 0, len(v.Skus)),
	}
	for _, sk := range v.Skus {
		res.Skus = append(res.Skus, v1.AdminSkuAdminItem{
			SkuId: fmtID(sk.SkuId), SkuNo: sk.SkuNo, Specs: sk.Specs, Price: sk.Price,
			LinePrice: sk.LinePrice, CostPrice: sk.CostPrice, Weight: sk.Weight,
			Barcode: sk.Barcode, Status: sk.Status,
		})
	}
	return res, nil
}
