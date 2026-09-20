package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// PayNotify 支付结果回调（开放路径; 渠道→平台, 需签名验证）
func (c *ControllerV1) PayNotify(ctx context.Context, req *v1.PayNotifyReq) (res *v1.PayNotifyRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
