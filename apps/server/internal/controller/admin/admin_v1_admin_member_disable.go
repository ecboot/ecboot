package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminMemberDisable 禁用/启用（审计留痕）
func (c *ControllerV1) AdminMemberDisable(ctx context.Context, req *v1.AdminMemberDisableReq) (res *v1.AdminMemberDisableRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
