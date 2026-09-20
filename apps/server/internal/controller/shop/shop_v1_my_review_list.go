package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// MyReviewList 我的评价
func (c *ControllerV1) MyReviewList(ctx context.Context, req *v1.MyReviewListReq) (res *v1.MyReviewListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
