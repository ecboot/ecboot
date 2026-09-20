package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminUserList 后台账号列表
func (c *ControllerV1) AdminUserList(ctx context.Context, req *v1.AdminUserListReq) (res *v1.AdminUserListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
