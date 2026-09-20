package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// CouponReceive 领券（防超发/限领, 幂等语义由服务端保证）
func (c *ControllerV1) CouponReceive(ctx context.Context, req *v1.CouponReceiveReq) (res *v1.CouponReceiveRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
