package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminSpuList SPU 列表（全状态）
func (c *ControllerV1) AdminSpuList(ctx context.Context, req *v1.AdminSpuListReq) (res *v1.AdminSpuListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
