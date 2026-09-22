package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/consts"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// AdminAfterSaleConfirmReceipt 退货确认收货（→待退款并立即发起退款）
// 权限: aftersale:audit
func (c *ControllerV1) AdminAfterSaleConfirmReceipt(ctx context.Context, req *v1.AdminAfterSaleConfirmReceiptReq) (res *v1.AdminAfterSaleConfirmReceiptRes, err error) {
	if err = middleware.RequirePerm(ctx, consts.PermAfterSaleAudit); err != nil {
		return nil, err
	}
	if err = shop.NewAfterSaleLogic().ConfirmReceipt(ctx, req.AfterSaleNo, adminOperator(ctx)); err != nil {
		return nil, err
	}
	return &v1.AdminAfterSaleConfirmReceiptRes{Success: true}, nil
}
