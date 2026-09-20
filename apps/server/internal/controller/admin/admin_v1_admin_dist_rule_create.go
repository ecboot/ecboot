package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistRuleCreate 创建佣金规则（作用域唯一）
func (c *ControllerV1) AdminDistRuleCreate(ctx context.Context, req *v1.AdminDistRuleCreateReq) (res *v1.AdminDistRuleCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
