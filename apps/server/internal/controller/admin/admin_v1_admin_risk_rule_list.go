package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRiskRuleList 风控规则列表
func (c *ControllerV1) AdminRiskRuleList(ctx context.Context, req *v1.AdminRiskRuleListReq) (res *v1.AdminRiskRuleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
