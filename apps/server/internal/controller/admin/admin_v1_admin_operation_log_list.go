package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminOperationLogList 操作审计日志
func (c *ControllerV1) AdminOperationLogList(ctx context.Context, req *v1.AdminOperationLogListReq) (res *v1.AdminOperationLogListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
