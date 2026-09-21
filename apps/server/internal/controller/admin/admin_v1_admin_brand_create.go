package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminBrandCreate 新增品牌
func (c *ControllerV1) AdminBrandCreate(ctx context.Context, req *v1.AdminBrandCreateReq) (res *v1.AdminBrandCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:brand:create"); err != nil {
		return nil, err
	}
	id, err := shop.NewProductLogic().AdminBrandCreate(ctx, model.BrandInput{
		Name: req.Name, Logo: req.Logo, Description: req.Description, Sort: req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminBrandCreateRes{Id: fmtID(id)}, nil
}
