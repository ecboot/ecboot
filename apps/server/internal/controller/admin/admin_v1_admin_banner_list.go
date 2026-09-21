package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminBannerList 轮播列表
func (c *ControllerV1) AdminBannerList(ctx context.Context, req *v1.AdminBannerListReq) (res *v1.AdminBannerListRes, err error) {
	out, err := shop.BannerList(ctx, req.Position, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminBannerListRes{List: make([]v1.AdminBannerItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminBannerItem{
			Id: fmtID(it.Id), Position: it.Position, ImageUrl: it.ImageUrl, LinkUrl: it.LinkUrl,
			Sort: it.Sort, StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status,
		})
	}
	return res, nil
}
