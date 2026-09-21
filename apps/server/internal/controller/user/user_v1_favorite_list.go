package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// FavoriteList 收藏列表（含实时价态与失效标注）
func (c *ControllerV1) FavoriteList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error) {
	out, err := user.FavoriteList(ctx, middleware.CtxUserIdFrom(ctx),
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.FavoriteListRes{List: make([]v1.FavoriteItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.FavoriteItem{
			SpuId: fmtID(it.SpuId), SpuName: it.Name, Image: it.Image,
			Price: it.Price, Sellable: it.Sellable, Invalid: it.Invalid,
		})
	}
	return res, nil
}
