package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAfterSaleRetryRefund 退款重试
func (c *ControllerV1) AdminAfterSaleRetryRefund(ctx context.Context, req *v1.AdminAfterSaleRetryRefundReq) (res *v1.AdminAfterSaleRetryRefundRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
