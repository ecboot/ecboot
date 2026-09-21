package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCategoryCreate 新增分类
func (c *ControllerV1) AdminCategoryCreate(ctx context.Context, req *v1.AdminCategoryCreateReq) (res *v1.AdminCategoryCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "product:category:create"); err != nil {
		return nil, err
	}
	var parentId int64
	if req.ParentId != "" && req.ParentId != "0" {
		if parentId, err = parseID(req.ParentId); err != nil {
			return nil, err
		}
	}
	id, err := shop.NewProductLogic().AdminCategoryCreate(ctx, model.CategoryInput{
		ParentId: parentId, Name: req.Name, Icon: req.Icon, Level: req.Level, Sort: req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminCategoryCreateRes{Id: fmtID(id)}, nil
}
