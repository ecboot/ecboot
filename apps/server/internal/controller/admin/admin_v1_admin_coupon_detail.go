package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminCouponDetail 券模板详情。
func (c *ControllerV1) AdminCouponDetail(ctx context.Context, req *v1.AdminCouponDetailReq) (res *v1.AdminCouponDetailRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermPromotionCouponRead); err != nil {
		return nil, err
	}
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	it, err := shop.NewCouponLogic().AdminDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.AdminCouponDetailRes{
		Id: fmtID(it.Id), Name: it.Name, Type: it.Type,
		Threshold: it.Threshold, Discount: it.Discount,
		TotalCount: it.TotalCount, PerLimit: it.PerLimit,
		ValidType: it.ValidType, ValidDesc: it.ValidDesc, Status: it.Status,
	}, nil
}
