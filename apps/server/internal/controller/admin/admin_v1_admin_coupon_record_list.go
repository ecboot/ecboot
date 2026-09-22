package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCouponRecordList 券领取/使用记录分页。
func (c *ControllerV1) AdminCouponRecordList(ctx context.Context, req *v1.AdminCouponRecordListReq) (res *v1.AdminCouponRecordListRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponRead); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewCouponLogic().AdminRecords(ctx, id, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminCouponRecordListRes{List: make([]v1.AdminCouponRecordItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.AdminCouponRecordItem{
			UserCouponId: fmtID(it.UserCouponId), UserId: fmtID(it.UserId),
			Status: it.Status, OrderNo: it.OrderNo, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
