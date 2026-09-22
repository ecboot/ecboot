package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCouponUpdate 修改券模板（status 显式启停: 0停发 1启用, 契约 dc 注明必传）。
func (c *ControllerV1) AdminCouponUpdate(ctx context.Context, req *v1.AdminCouponUpdateReq) (res *v1.AdminCouponUpdateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponUpdate); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = shop.NewCouponLogic().AdminUpdate(ctx, id, model.CouponInput{
		Name: req.Name, Threshold: req.Threshold, Discount: req.Discount,
		TotalCount: req.TotalCount, PerLimit: req.PerLimit, Status: req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.AdminCouponUpdateRes{Success: true}, nil
}
