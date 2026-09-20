package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminBannerList 轮播管理
func (c *ControllerV1) AdminBannerList(ctx context.Context, req *v1.AdminBannerListReq) (res *v1.AdminBannerListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
