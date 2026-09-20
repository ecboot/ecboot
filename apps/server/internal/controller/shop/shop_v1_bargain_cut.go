package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BargainCut 帮砍（会员; 一人一刀; 风控挂载位: risk_rule 2/3/4）
func (c *ControllerV1) BargainCut(ctx context.Context, req *v1.BargainCutReq) (res *v1.BargainCutRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
