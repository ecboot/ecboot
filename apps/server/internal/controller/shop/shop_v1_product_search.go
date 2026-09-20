package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// ProductSearch 关键词搜索
func (c *ControllerV1) ProductSearch(ctx context.Context, req *v1.ProductSearchReq) (res *v1.ProductSearchRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
