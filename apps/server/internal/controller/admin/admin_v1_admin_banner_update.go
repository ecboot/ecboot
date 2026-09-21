package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminBannerUpdate 修改轮播（全量覆盖: status 必传——空即视为停用）
func (c *ControllerV1) AdminBannerUpdate(ctx context.Context, req *v1.AdminBannerUpdateReq) (res *v1.AdminBannerUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, "operation:banner:manage"); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	// position 传 0: api 契约无该字段, 位置创建后不可改（service 不改写 position 列）
	if err = shop.BannerUpdate(ctx, id, bannerInputFromReq(
		0, req.ImageUrl, req.LinkUrl, req.Sort, req.StartTime, req.EndTime, req.Status,
	)); err != nil {
		return nil, err
	}
	return &v1.AdminBannerUpdateRes{Success: true}, nil
}
