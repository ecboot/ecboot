package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAfterSaleConfirmReceipt 退货确认收货（→待退款）
func (c *ControllerV1) AdminAfterSaleConfirmReceipt(ctx context.Context, req *v1.AdminAfterSaleConfirmReceiptReq) (res *v1.AdminAfterSaleConfirmReceiptRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
