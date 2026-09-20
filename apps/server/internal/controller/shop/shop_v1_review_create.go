package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// ReviewCreate 提交评价（一项一评; 004 契约, 复用 V12 表）
func (c *ControllerV1) ReviewCreate(ctx context.Context, req *v1.ReviewCreateReq) (res *v1.ReviewCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
