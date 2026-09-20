package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminCategoryTree 分类树（含禁用）
func (c *ControllerV1) AdminCategoryTree(ctx context.Context, req *v1.AdminCategoryTreeReq) (res *v1.AdminCategoryTreeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
