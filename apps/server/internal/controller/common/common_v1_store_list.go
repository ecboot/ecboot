package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

// StoreList 门店列表（区县筛选或经纬度附近检索, 距离排序）—— contracts/common-api.md
func (c *ControllerV1) StoreList(ctx context.Context, req *v1.StoreListReq) (res *v1.StoreListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
