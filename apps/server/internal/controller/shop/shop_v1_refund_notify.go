package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// RefundNotify 退款结果回调（开放路径; 需签名验证）
func (c *ControllerV1) RefundNotify(ctx context.Context, req *v1.RefundNotifyReq) (res *v1.RefundNotifyRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
