package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

func (c *ControllerV1) AdminRiskRuleDelete(ctx context.Context, req *v1.AdminRiskRuleDeleteReq) (res *v1.AdminRiskRuleDeleteRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
