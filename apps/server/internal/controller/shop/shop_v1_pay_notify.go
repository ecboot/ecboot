package shop

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// PayNotify 支付结果回调（渠道; 公开端点——签名由渠道适配层校验）
func (c *ControllerV1) PayNotify(ctx context.Context, req *v1.PayNotifyReq) (res *v1.PayNotifyRes, err error) {
	raw := g.RequestFromCtx(ctx).GetBody()
	r := g.RequestFromCtx(ctx)
	if err = shop.NewPayLogic().HandlePayNotify(ctx, "mock", raw); err != nil {
		// 评审 I10: 渠道应答须**直写响应体**（照 middleware.unauthorized 先例, Response 中间件不覆盖）——
		// 否则渠道收到的是统一三段式包装（恒 code:0）, "失败→重试"语义不成立; 且不泄露内部错误文案。
		g.Log().Errorf(ctx, "支付回调处理失败: %v", err)
		r.Response.WriteJsonExit(g.Map{"code": "FAIL", "message": "处理失败"})
		return nil, nil
	}
	r.Response.WriteJsonExit(g.Map{"code": "SUCCESS", "message": "OK"})
	return nil, nil
}
