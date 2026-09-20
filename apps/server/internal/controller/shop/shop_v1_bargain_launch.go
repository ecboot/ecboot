package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BargainLaunch 发起砍价（会员; 选择场次商品开一刀）
func (c *ControllerV1) BargainLaunch(ctx context.Context, req *v1.BargainLaunchReq) (res *v1.BargainLaunchRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
