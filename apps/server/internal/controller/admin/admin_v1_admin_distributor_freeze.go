package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistributorFreeze 冻结/解冻
func (c *ControllerV1) AdminDistributorFreeze(ctx context.Context, req *v1.AdminDistributorFreezeReq) (res *v1.AdminDistributorFreezeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
