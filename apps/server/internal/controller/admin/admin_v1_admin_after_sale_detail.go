package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleDetail 售后详情
// 权限: aftersale:read
func (c *ControllerV1) AdminAfterSaleDetail(ctx context.Context, req *v1.AdminAfterSaleDetailReq) (res *v1.AdminAfterSaleDetailRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleRead); err != nil {
		return nil, err
	}
	d, err := shop.NewAfterSaleLogic().AdminDetail(ctx, req.AfterSaleNo)
	if err != nil {
		return nil, err
	}
	return &v1.AdminAfterSaleDetailRes{
		AfterSaleNo:       d.AfterSaleNo,
		OrderNo:           d.OrderNo,
		UserId:            fmtID(d.UserId),
		Type:              d.Type,
		Quantity:          d.Quantity,
		Reason:            d.Reason,
		Description:       d.Description,
		VoucherImages:     d.VoucherImages,
		RefundAmount:      d.RefundAmount,
		ReturnLogisticsNo: d.ReturnLogisticsNo,
		RejectReason:      d.RejectReason,
		RefundNo:          d.RefundNo,
		Status:            d.Status,
	}, nil
}
