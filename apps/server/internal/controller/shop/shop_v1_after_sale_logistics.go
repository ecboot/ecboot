package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

func (c *ControllerV1) AfterSaleLogistics(ctx context.Context, req *v1.AfterSaleLogisticsReq) (res *v1.AfterSaleLogisticsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
