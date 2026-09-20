package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminRiskRecordList 风控事件查询
func (c *ControllerV1) AdminRiskRecordList(ctx context.Context, req *v1.AdminRiskRecordListReq) (res *v1.AdminRiskRecordListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
