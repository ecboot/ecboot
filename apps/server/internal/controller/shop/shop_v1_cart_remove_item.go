package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

func (c *ControllerV1) CartRemoveItem(ctx context.Context, req *v1.CartRemoveItemReq) (res *v1.CartRemoveItemRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
