package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminWithdrawList 提现列表
func (c *ControllerV1) AdminWithdrawList(ctx context.Context, req *v1.AdminWithdrawListReq) (res *v1.AdminWithdrawListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
