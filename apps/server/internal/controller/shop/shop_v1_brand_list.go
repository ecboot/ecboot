package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BrandList 品牌列表（启用）
func (c *ControllerV1) BrandList(ctx context.Context, req *v1.BrandListReq) (res *v1.BrandListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
