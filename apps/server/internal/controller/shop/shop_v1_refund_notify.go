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
	r := g.RequestFromCtx(ctx)
	if err = shop.NewPayLogic().HandleRefundNotify(ctx, "mock", raw); err != nil {
		// 评审 I10: 同支付回调——直写响应体, 不泄露内部文案
		g.Log().Errorf(ctx, "退款回调处理失败: %v", err)
		r.Response.WriteJsonExit(g.Map{"code": "FAIL", "message": "处理失败"})
		return nil, nil
	}
	r.Response.WriteJsonExit(g.Map{"code": "SUCCESS", "message": "OK"})
	return nil, nil
}
