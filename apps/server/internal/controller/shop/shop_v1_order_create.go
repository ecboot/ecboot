package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// OrderCreate 创建订单（幂等; 玩法上下文二选一见字段）
func (c *ControllerV1) OrderCreate(ctx context.Context, req *v1.OrderCreateReq) (res *v1.OrderCreateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
