package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminInventoryWarn 预警列表（可售≤阈值）
func (c *ControllerV1) AdminInventoryWarn(ctx context.Context, req *v1.AdminInventoryWarnReq) (res *v1.AdminInventoryWarnRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
