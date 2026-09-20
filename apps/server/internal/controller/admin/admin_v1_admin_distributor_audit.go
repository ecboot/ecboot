package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistributorAudit 推广员审核
func (c *ControllerV1) AdminDistributorAudit(ctx context.Context, req *v1.AdminDistributorAuditReq) (res *v1.AdminDistributorAuditRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
