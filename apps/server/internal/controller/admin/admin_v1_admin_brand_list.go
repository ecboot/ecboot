package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminBrandList 品牌列表
func (c *ControllerV1) AdminBrandList(ctx context.Context, req *v1.AdminBrandListReq) (res *v1.AdminBrandListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
