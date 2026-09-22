package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminAssistList 助力活动列表。
func (c *ControllerV1) AdminAssistList(ctx context.Context, req *v1.AdminAssistListReq) (res *v1.AdminAssistListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionAssistAll); err != nil {
		return nil, err
	}
	out, err := shop.NewActivityLogic().AssistList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminAssistListRes{List: make([]v1.AdminAssistItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminAssistItem{Id: fmtID(it.Id), Name: it.Name, RewardType: it.RewardType, RewardDesc: it.RewardDesc, RequiredCount: it.RequiredCount, PerLimit: it.PerLimit, StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status})
	}
	return res, nil
}
