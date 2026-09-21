package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminFloorDelete 删除楼层(软删)
func (c *ControllerV1) AdminFloorDelete(ctx context.Context, req *v1.AdminFloorDeleteReq) (res *v1.AdminFloorDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:floor:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.FloorDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminFloorDeleteRes{Success: true}, nil
}
