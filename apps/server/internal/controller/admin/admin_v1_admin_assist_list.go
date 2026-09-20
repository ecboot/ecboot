package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminAssistList ---------- 助力 ----------
func (c *ControllerV1) AdminAssistList(ctx context.Context, req *v1.AdminAssistListReq) (res *v1.AdminAssistListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
