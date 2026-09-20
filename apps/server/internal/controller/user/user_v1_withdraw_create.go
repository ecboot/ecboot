package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

// WithdrawCreate 提现申请
func (c *ControllerV1) WithdrawCreate(ctx context.Context, req *v1.WithdrawCreateReq) (res *v1.WithdrawCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
