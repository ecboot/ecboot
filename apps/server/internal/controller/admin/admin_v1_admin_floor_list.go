package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminFloorList 楼层管理
func (c *ControllerV1) AdminFloorList(ctx context.Context, req *v1.AdminFloorListReq) (res *v1.AdminFloorListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
