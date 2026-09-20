package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// FloorList 楼层内容（含商品摘要装配）
func (c *ControllerV1) FloorList(ctx context.Context, req *v1.FloorListReq) (res *v1.FloorListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
