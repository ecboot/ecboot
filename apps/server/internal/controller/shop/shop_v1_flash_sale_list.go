package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// FlashSaleList 秒杀列表（进行中与预告）
func (c *ControllerV1) FlashSaleList(ctx context.Context, req *v1.FlashSaleListReq) (res *v1.FlashSaleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
