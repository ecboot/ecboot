package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminInviteRecordList 邀请激励记录
func (c *ControllerV1) AdminInviteRecordList(ctx context.Context, req *v1.AdminInviteRecordListReq) (res *v1.AdminInviteRecordListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
