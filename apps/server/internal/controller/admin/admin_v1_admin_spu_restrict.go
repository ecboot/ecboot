package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminSpuRestrict 限售区域设置
func (c *ControllerV1) AdminSpuRestrict(ctx context.Context, req *v1.AdminSpuRestrictReq) (res *v1.AdminSpuRestrictRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
