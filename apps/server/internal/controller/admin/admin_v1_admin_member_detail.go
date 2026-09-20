package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminMemberDetail 会员详情（含资产概要）
func (c *ControllerV1) AdminMemberDetail(ctx context.Context, req *v1.AdminMemberDetailReq) (res *v1.AdminMemberDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
