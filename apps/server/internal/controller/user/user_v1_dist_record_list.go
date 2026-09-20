package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
)

func (c *ControllerV1) DistRecordList(ctx context.Context, req *v1.DistRecordListReq) (res *v1.DistRecordListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
