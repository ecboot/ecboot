package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminSpuCreate 创建商品SPU
func (c *ControllerV1) AdminSpuCreate(ctx context.Context, req *v1.AdminSpuCreateReq) (res *v1.AdminSpuCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:spu:create"); err != nil {
		return nil, err
	}
	catId, err := parseID(req.CategoryId)
	if err != nil {
		return nil, err
	}
	brandId, ftId, err := optIDPair(req.BrandId, req.FreightTemplateId)
	if err != nil {
		return nil, err
	}
	id, spuNo, err := shop.NewProductLogic().AdminProductCreate(ctx, model.SpuInput{
		Name: req.Name, SubTitle: req.SubTitle, CategoryId: catId, BrandId: brandId,
		FreightTemplateId: ftId, Images: req.Images, VideoUrl: req.VideoUrl,
		Description: req.Description, SpecDefinitions: req.SpecDefinitions, Attributes: req.Attributes,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminSpuCreateRes{SpuId: fmtID(id), SpuNo: spuNo}, nil
}
