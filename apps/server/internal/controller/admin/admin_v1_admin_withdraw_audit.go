package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminWithdrawAudit 提现审核（通过→打款中冻结; 拒绝→回退）
func (c *ControllerV1) AdminWithdrawAudit(ctx context.Context, req *v1.AdminWithdrawAuditReq) (res *v1.AdminWithdrawAuditRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
