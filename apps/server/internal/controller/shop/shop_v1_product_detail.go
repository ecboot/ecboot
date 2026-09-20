package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// ProductDetail 商品详情（SKU/规格/运费概要/评价汇总; 不含成本价）
func (c *ControllerV1) ProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
