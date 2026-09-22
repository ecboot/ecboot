package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminCouponCreate 创建券模板。
func (c *ControllerV1) AdminCouponCreate(ctx context.Context, req *v1.AdminCouponCreateReq) (res *v1.AdminCouponCreateRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponCreate); err != nil {
		return nil, err
	}
	id, err := shop.NewCouponLogic().AdminCreate(ctx, model.CouponInput{
		Name: req.Name, Type: req.Type, Threshold: req.Threshold, Discount: req.Discount,
		TotalCount: req.TotalCount, PerLimit: req.PerLimit,
		ValidType: req.ValidType, ValidStartAt: req.ValidStartAt, ValidEndAt: req.ValidEndAt, ValidDays: req.ValidDays,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminCouponCreateRes{Id: fmtID(id)}, nil
}
