package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// ReviewExtra 追评（一次, 90 天内）
func (c *ControllerV1) ReviewExtra(ctx context.Context, req *v1.ReviewExtraReq) (res *v1.ReviewExtraRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
