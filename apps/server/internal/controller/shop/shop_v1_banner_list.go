package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BannerList 轮播/弹窗（投放时段内, 按位置）
func (c *ControllerV1) BannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
