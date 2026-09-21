package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// CouponAvailableList 可领优惠券列表（含个人可领标记）
func (c *ControllerV1) CouponAvailableList(ctx context.Context, req *v1.CouponAvailableListReq) (res *v1.CouponAvailableListRes, err error) {
	out, err := user.AvailableTemplates(ctx, middleware.CtxUserIdFrom(ctx),
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.CouponAvailableListRes{List: make([]v1.UsableCouponTemplate, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.UsableCouponTemplate{
			CouponId: fmtID(it.CouponId), Name: it.Name, Type: it.Type,
			Threshold: it.Threshold, Discount: it.Discount,
			ValidDesc: it.ValidDesc, CanReceive: it.CanReceive,
		})
	}
	return res, nil
}
