package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRoleDetail 角色详情
func (c *ControllerV1) AdminRoleDetail(ctx context.Context, req *v1.AdminRoleDetailReq) (res *v1.AdminRoleDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
