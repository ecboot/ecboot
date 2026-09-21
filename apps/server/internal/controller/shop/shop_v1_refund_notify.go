package shop

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// RefundNotify 退款结果回调（渠道; 公开端点）
func (c *ControllerV1) RefundNotify(ctx context.Context, req *v1.RefundNotifyReq) (res *v1.RefundNotifyRes, err error) {
	raw := g.RequestFromCtx(ctx).GetBody()
	if err = shop.NewPayLogic().HandleRefundNotify(ctx, "mock", raw); err != nil {
		return &v1.RefundNotifyRes{Code: "FAIL", Message: err.Error()}, nil
	}
	return &v1.RefundNotifyRes{Code: "SUCCESS", Message: "OK"}, nil
}
