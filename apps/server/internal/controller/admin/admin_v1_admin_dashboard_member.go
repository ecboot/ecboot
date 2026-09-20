package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDashboardMember 会员看板
func (c *ControllerV1) AdminDashboardMember(ctx context.Context, req *v1.AdminDashboardMemberReq) (res *v1.AdminDashboardMemberRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
