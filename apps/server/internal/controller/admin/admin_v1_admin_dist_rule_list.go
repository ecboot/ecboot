package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistRuleList 佣金规则列表
func (c *ControllerV1) AdminDistRuleList(ctx context.Context, req *v1.AdminDistRuleListReq) (res *v1.AdminDistRuleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
