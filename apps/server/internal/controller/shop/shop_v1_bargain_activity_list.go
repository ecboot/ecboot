package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BargainActivityList 砍价活动列表
func (c *ControllerV1) BargainActivityList(ctx context.Context, req *v1.BargainActivityListReq) (res *v1.BargainActivityListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
