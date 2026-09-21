package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminBrandUpdate 修改品牌
func (c *ControllerV1) AdminBrandUpdate(ctx context.Context, req *v1.AdminBrandUpdateReq) (res *v1.AdminBrandUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:brand:update"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminBrandUpdate(ctx, id, model.BrandInput{
		Name: req.Name, Logo: req.Logo, Description: req.Description,
		Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminBrandUpdateRes{Success: true}, nil
}
