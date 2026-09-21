package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminLogisticsDelete 删除物流公司(软删)
func (c *ControllerV1) AdminLogisticsDelete(ctx context.Context, req *v1.AdminLogisticsDeleteReq) (res *v1.AdminLogisticsDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "logistics:company:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.LogisticsDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminLogisticsDeleteRes{Success: true}, nil
}
