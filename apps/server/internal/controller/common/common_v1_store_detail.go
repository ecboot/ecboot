package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

// StoreDetail 门店详情
func (c *ControllerV1) StoreDetail(ctx context.Context, req *v1.StoreDetailReq) (res *v1.StoreDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
