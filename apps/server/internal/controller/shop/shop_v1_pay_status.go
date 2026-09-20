package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// PayStatus 支付状态查询
func (c *ControllerV1) PayStatus(ctx context.Context, req *v1.PayStatusReq) (res *v1.PayStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
