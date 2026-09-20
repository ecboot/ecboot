package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminLogisticsDetail 物流公司详情
func (c *ControllerV1) AdminLogisticsDetail(ctx context.Context, req *v1.AdminLogisticsDetailReq) (res *v1.AdminLogisticsDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
