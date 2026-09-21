package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminBannerDelete 删除轮播(软删)
func (c *ControllerV1) AdminBannerDelete(ctx context.Context, req *v1.AdminBannerDeleteReq) (res *v1.AdminBannerDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:banner:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.BannerDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminBannerDeleteRes{Success: true}, nil
}
