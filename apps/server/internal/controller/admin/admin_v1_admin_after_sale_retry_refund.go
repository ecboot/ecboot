package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleRetryRefund 退款重试（渠道失败后; 仅"待退款"可重试）
// 权限: aftersale:refund
func (c *ControllerV1) AdminAfterSaleRetryRefund(ctx context.Context, req *v1.AdminAfterSaleRetryRefundReq) (res *v1.AdminAfterSaleRetryRefundRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleRefund); err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().RetryRefund(ctx, req.AfterSaleNo, adminOperator(ctx)); err != nil {
		return nil, err
	}
	return &v1.AdminAfterSaleRetryRefundRes{Success: true}, nil
}
