package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminBannerCreate 新增轮播（新建默认启用）
func (c *ControllerV1) AdminBannerCreate(ctx context.Context, req *v1.AdminBannerCreateReq) (res *v1.AdminBannerCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:banner:manage"); err != nil {
		return nil, err
	}
	id, err := shop.BannerCreate(ctx, bannerInputFromReq(
		req.Position, req.ImageUrl, req.LinkUrl, req.Sort, req.StartTime, req.EndTime, 1))
	if err != nil {
		return nil, err
	}
	return &v1.AdminBannerCreateRes{Id: fmtID(id)}, nil
}
