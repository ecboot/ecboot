package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminSpuCreate 创建 SPU（含规格定义/图文/运费模板）
func (c *ControllerV1) AdminSpuCreate(ctx context.Context, req *v1.AdminSpuCreateReq) (res *v1.AdminSpuCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
