package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)


// AdminFlashSaleList 秒杀活动列表。
func (c *ControllerV1) AdminFlashSaleList(ctx context.Context, req *v1.AdminFlashSaleListReq) (res *v1.AdminFlashSaleListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionFlashSaleAll); err != nil {
		return nil, err
	}
	out, err := shop.NewActivityLogic().FlashSaleList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminFlashSaleListRes{List: make([]v1.AdminFlashSaleItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminFlashSaleItem{Id: fmtID(it.Id), Name: it.Name, StartTime: it.StartTime, EndTime: it.EndTime, Status: it.Status})
	}
	return res, nil
}
