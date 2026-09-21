package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminFloorUpdate 修改楼层（类型不可改——api 契约无 floorType; status 必传）
func (c *ControllerV1) AdminFloorUpdate(ctx context.Context, req *v1.AdminFloorUpdateReq) (res *v1.AdminFloorUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:floor:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.FloorUpdate(ctx, id, floorInputFromReq(0, req.Title, req.Config, req.Sort, req.Status)); err != nil {
		return nil, err
	}
	return &v1.AdminFloorUpdateRes{Success: true}, nil
}
