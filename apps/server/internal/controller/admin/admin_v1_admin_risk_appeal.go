package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRiskAppeal 申诉处理
func (c *ControllerV1) AdminRiskAppeal(ctx context.Context, req *v1.AdminRiskAppealReq) (res *v1.AdminRiskAppealRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
