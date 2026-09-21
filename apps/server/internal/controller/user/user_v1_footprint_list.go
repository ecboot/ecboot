package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// FootprintList 浏览足迹（最近浏览倒序）
func (c *ControllerV1) FootprintList(ctx context.Context, req *v1.FootprintListReq) (res *v1.FootprintListRes, err error) {
	out, err := user.FootprintList(ctx, middleware.CtxUserIdFrom(ctx),
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.FootprintListRes{List: make([]v1.FootprintItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.FootprintItem{
			SpuId: fmtID(it.SpuId), SpuName: it.Name, Image: it.Image,
			Price: it.Price, LastViewAt: it.LastViewAt,
		})
	}
	return res, nil
}
