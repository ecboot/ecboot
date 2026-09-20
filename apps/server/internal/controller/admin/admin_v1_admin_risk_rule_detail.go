package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRiskRuleDetail 风控规则详情
func (c *ControllerV1) AdminRiskRuleDetail(ctx context.Context, req *v1.AdminRiskRuleDetailReq) (res *v1.AdminRiskRuleDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
