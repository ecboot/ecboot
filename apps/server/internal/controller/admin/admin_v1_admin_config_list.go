package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminConfigList 系统配置列表
func (c *ControllerV1) AdminConfigList(ctx context.Context, req *v1.AdminConfigListReq) (res *v1.AdminConfigListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
