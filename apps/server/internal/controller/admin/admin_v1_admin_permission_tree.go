package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminPermissionTree 权限树（菜单/按钮/接口统一树）
func (c *ControllerV1) AdminPermissionTree(ctx context.Context, req *v1.AdminPermissionTreeReq) (res *v1.AdminPermissionTreeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
