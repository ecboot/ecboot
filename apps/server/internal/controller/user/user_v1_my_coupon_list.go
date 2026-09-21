package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// MyCouponList 我的优惠券（四态; 过期惰性判定）
func (c *ControllerV1) MyCouponList(ctx context.Context, req *v1.MyCouponListReq) (res *v1.MyCouponListRes, err error) {
	out, err := user.Mine(ctx, middleware.CtxUserIdFrom(ctx), req.Status,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.MyCouponListRes{List: make([]v1.MyCouponItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.MyCouponItem{
			UserCouponId: fmtID(it.UserCouponId), Name: it.Name,
			Threshold: it.Threshold, Discount: it.Discount,
			ExpireTime: it.ExpireTime, Status: it.Status,
		})
	}
	return res, nil
}
