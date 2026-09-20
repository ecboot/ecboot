package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// CartCheckout 结算试算（勾选项: 商品价态+优惠明细+运费+应付）
func (c *ControllerV1) CartCheckout(ctx context.Context, req *v1.CartCheckoutReq) (res *v1.CartCheckoutRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
