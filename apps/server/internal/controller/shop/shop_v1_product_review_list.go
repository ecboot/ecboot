package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// ProductReviewList 商品评价列表（审核通过）
func (c *ControllerV1) ProductReviewList(ctx context.Context, req *v1.ProductReviewListReq) (res *v1.ProductReviewListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
