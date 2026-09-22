package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCouponList 券模板列表（status 筛选 + 分页; 含已领数）。
func (c *ControllerV1) AdminCouponList(ctx context.Context, req *v1.AdminCouponListReq) (res *v1.AdminCouponListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponRead); err != nil {
		return nil, err
	}
	out, err := shop.NewCouponLogic().AdminList(ctx, req.Status, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminCouponListRes{List: make([]v1.AdminCouponItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminCouponItem{
			Id: fmtID(it.Id), Name: it.Name, Type: it.Type,
			Threshold: it.Threshold, Discount: it.Discount,
			TotalCount: it.TotalCount, ReceivedCount: it.Received, PerLimit: it.PerLimit,
			ValidType: it.ValidType, ValidDesc: it.ValidDesc, Status: it.Status,
		})
	}
	return res, nil
}
