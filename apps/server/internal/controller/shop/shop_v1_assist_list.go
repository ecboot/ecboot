package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// AssistList 助力活动列表
func (c *ControllerV1) AssistList(ctx context.Context, req *v1.AssistListReq) (res *v1.AssistListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
