package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminSpuUpdate 修改商品SPU
func (c *ControllerV1) AdminSpuUpdate(ctx context.Context, req *v1.AdminSpuUpdateReq) (res *v1.AdminSpuUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:spu:update"); err != nil {
		return nil, err
	}
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	catId, brandId, ftId, err := optIDTriple(req.CategoryId, req.BrandId, req.FreightTemplateId)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminProductUpdate(ctx, spuId, model.SpuInput{
		Name: req.Name, SubTitle: req.SubTitle, CategoryId: catId, BrandId: brandId,
		FreightTemplateId: ftId, Images: req.Images, VideoUrl: req.VideoUrl,
		Description: req.Description, SpecDefinitions: req.SpecDefinitions, Attributes: req.Attributes,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminSpuUpdateRes{Success: true}, nil
}
