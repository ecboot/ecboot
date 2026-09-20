package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminWithdrawPay 打款结果登记（成功核销 / 失败回退）
func (c *ControllerV1) AdminWithdrawPay(ctx context.Context, req *v1.AdminWithdrawPayReq) (res *v1.AdminWithdrawPayRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
