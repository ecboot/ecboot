package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// BannerList 运营位列表（轮播/弹窗; 在投过滤, 公开）
func (c *ControllerV1) BannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error) {
	list, err := shop.PublicBanners(ctx, req.Position)
	if err != nil {
		return nil, err
	}
	res = &v1.BannerListRes{List: make([]v1.BannerItem, 0, len(list))}
	for _, it := range list {
		res.List = append(res.List, v1.BannerItem{
			Id: fmtID(it.Id), ImageUrl: it.ImageUrl, LinkUrl: it.LinkUrl,
		})
	}
	return res, nil
}
