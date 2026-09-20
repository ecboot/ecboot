package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminInventoryList 库存列表（可售/锁定）
func (c *ControllerV1) AdminInventoryList(ctx context.Context, req *v1.AdminInventoryListReq) (res *v1.AdminInventoryListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
