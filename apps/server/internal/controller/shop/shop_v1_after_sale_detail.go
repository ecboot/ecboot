package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

func (c *ControllerV1) AfterSaleDetail(ctx context.Context, req *v1.AfterSaleDetailReq) (res *v1.AfterSaleDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
