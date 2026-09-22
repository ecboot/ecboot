package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminBargainList 砍价活动列表。
func (c *ControllerV1) AdminBargainList(ctx context.Context, req *v1.AdminBargainListReq) (res *v1.AdminBargainListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionBargainAll); err != nil {
		return nil, err
	}
	out, err := shop.NewActivityLogic().BargainList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminBargainListRes{List: make([]v1.AdminBargainItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminBargainItem{Id: fmtID(it.Id), Name: it.Name, SpuId: fmtID(it.SpuId), StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status})
	}
	return res, nil
}
