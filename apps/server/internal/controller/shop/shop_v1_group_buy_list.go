package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// GroupBuyList 拼团活动列表（进行中）
func (c *ControllerV1) GroupBuyList(ctx context.Context, req *v1.GroupBuyListReq) (res *v1.GroupBuyListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
