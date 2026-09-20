package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// FullReductionList 当前满减活动（可按商品过滤范围命中）
func (c *ControllerV1) FullReductionList(ctx context.Context, req *v1.FullReductionListReq) (res *v1.FullReductionListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
