package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminDistRecordList 佣金记录（全局）
func (c *ControllerV1) AdminDistRecordList(ctx context.Context, req *v1.AdminDistRecordListReq) (res *v1.AdminDistRecordListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
