package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// PayCreate 发起支付（创建支付单; 返回渠道唤起参数）
func (c *ControllerV1) PayCreate(ctx context.Context, req *v1.PayCreateReq) (res *v1.PayCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
