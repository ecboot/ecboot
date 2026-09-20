package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// AfterSaleCreate 申请售后（按订单项; 仅退款/退货退款）
func (c *ControllerV1) AfterSaleCreate(ctx context.Context, req *v1.AfterSaleCreateReq) (res *v1.AfterSaleCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
