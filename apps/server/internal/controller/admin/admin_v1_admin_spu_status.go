package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminSpuStatus SPU 上下架（无启用 SKU 禁上架）
func (c *ControllerV1) AdminSpuStatus(ctx context.Context, req *v1.AdminSpuStatusReq) (res *v1.AdminSpuStatusRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
