package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAfterSaleApprove 同意（仅退款→待退款; 退货退款→待寄回）
func (c *ControllerV1) AdminAfterSaleApprove(ctx context.Context, req *v1.AdminAfterSaleApproveReq) (res *v1.AdminAfterSaleApproveRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
