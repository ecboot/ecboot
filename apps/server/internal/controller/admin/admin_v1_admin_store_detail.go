package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminStoreDetail 门店详情（管理）
func (c *ControllerV1) AdminStoreDetail(ctx context.Context, req *v1.AdminStoreDetailReq) (res *v1.AdminStoreDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
