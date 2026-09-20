package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminGroupBuyList ---------- 拼团 ----------
func (c *ControllerV1) AdminGroupBuyList(ctx context.Context, req *v1.AdminGroupBuyListReq) (res *v1.AdminGroupBuyListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
