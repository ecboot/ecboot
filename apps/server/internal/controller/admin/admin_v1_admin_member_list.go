package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminMemberList 会员列表（手机号精确检索; 脱敏展示）
func (c *ControllerV1) AdminMemberList(ctx context.Context, req *v1.AdminMemberListReq) (res *v1.AdminMemberListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
