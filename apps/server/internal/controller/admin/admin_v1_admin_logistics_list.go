package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminLogisticsList 物流公司字典列表
func (c *ControllerV1) AdminLogisticsList(ctx context.Context, req *v1.AdminLogisticsListReq) (res *v1.AdminLogisticsListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
