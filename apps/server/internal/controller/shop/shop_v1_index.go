package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
