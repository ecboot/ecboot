package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminFloorCreate 新增楼层（新建默认启用）
func (c *ControllerV1) AdminFloorCreate(ctx context.Context, req *v1.AdminFloorCreateReq) (res *v1.AdminFloorCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:floor:manage"); err != nil {
		return nil, err
	}
	id, err := shop.FloorCreate(ctx, floorInputFromReq(req.FloorType, req.Title, req.Config, req.Sort, 1))
	if err != nil {
		return nil, err
	}
	return &v1.AdminFloorCreateRes{Id: fmtID(id)}, nil
}
