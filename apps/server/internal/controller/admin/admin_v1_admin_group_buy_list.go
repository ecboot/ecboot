package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminGroupBuyList 拼团活动列表。
func (c *ControllerV1) AdminGroupBuyList(ctx context.Context, req *v1.AdminGroupBuyListReq) (res *v1.AdminGroupBuyListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionGroupBuyAll); err != nil {
		return nil, err
	}
	out, err := shop.NewActivityLogic().GroupBuyList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminGroupBuyListRes{List: make([]v1.AdminGroupBuyItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminGroupBuyItem{Id: fmtID(it.Id), Name: it.Name, SpuId: fmtID(it.SpuId), GroupSize: it.GroupSize, PerLimit: it.PerLimit, StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status})
	}
	return res, nil
}
