package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistributorList 推广员列表（审核/冻结入口）
func (c *ControllerV1) AdminDistributorList(ctx context.Context, req *v1.AdminDistributorListReq) (res *v1.AdminDistributorListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
