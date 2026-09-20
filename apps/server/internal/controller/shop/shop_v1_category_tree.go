package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// CategoryTree 三级分类树（禁用不返回）
func (c *ControllerV1) CategoryTree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
