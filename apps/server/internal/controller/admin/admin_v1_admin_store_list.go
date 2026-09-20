package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminStoreList 门店列表（管理）
func (c *ControllerV1) AdminStoreList(ctx context.Context, req *v1.AdminStoreListReq) (res *v1.AdminStoreListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
