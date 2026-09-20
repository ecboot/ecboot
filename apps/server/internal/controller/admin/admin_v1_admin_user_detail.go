package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminUserDetail 后台账号详情
func (c *ControllerV1) AdminUserDetail(ctx context.Context, req *v1.AdminUserDetailReq) (res *v1.AdminUserDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
