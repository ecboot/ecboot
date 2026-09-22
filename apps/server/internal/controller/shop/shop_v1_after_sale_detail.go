package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AfterSaleDetail 售后详情（仅本人; 他人按不存在）
func (c *ControllerV1) AfterSaleDetail(ctx context.Context, req *v1.AfterSaleDetailReq) (res *v1.AfterSaleDetailRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	d, err := shop.NewAfterSaleLogic().Detail(ctx, userId, req.AfterSaleNo)
	if err != nil {
		return nil, err
	}
	return &v1.AfterSaleDetailRes{
		AfterSaleNo:       d.AfterSaleNo,
		Status:            d.Status,
		Type:              d.Type,
		Quantity:          d.Quantity,
		Reason:            d.Reason,
		Description:       d.Description,
		VoucherImages:     d.VoucherImages,
		RefundAmount:      d.RefundAmount,
		ReturnLogisticsNo: d.ReturnLogisticsNo,
		RejectReason:      d.RejectReason,
		AuditTime:         d.AuditTime,
		RefundTime:        d.RefundTime,
	}, nil
}
