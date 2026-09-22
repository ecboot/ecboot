package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminCouponDelete 软删券模板（既有 user_coupon 不受影响）。
func (c *ControllerV1) AdminCouponDelete(ctx context.Context, req *v1.AdminCouponDeleteReq) (res *v1.AdminCouponDeleteRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponDelete); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewCouponLogic().AdminDelete(ctx, id); err != nil {
		return nil, err
	}
	return &v1.AdminCouponDeleteRes{Success: true}, nil
}
