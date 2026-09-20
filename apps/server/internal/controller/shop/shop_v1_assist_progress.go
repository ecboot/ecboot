package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// AssistProgress 助力进度（公开）
func (c *ControllerV1) AssistProgress(ctx context.Context, req *v1.AssistProgressReq) (res *v1.AssistProgressRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
