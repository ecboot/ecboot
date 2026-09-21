package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCategoryUpdate 修改分类
func (c *ControllerV1) AdminCategoryUpdate(ctx context.Context, req *v1.AdminCategoryUpdateReq) (res *v1.AdminCategoryUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:category:update"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewProductLogic().AdminCategoryUpdate(ctx, id, model.CategoryInput{
		Name: req.Name, Icon: req.Icon, Sort: req.Sort, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminCategoryUpdateRes{Success: true}, nil
}
