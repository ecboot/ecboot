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
	if err = shop.NewPayLogic().HandlePayNotify(ctx, "mock", raw); err != nil {
		// 渠道应答: 失败语义（渠道将重试）
		return &v1.PayNotifyRes{Code: "FAIL", Message: err.Error()}, nil
	}
	return &v1.PayNotifyRes{Code: "SUCCESS", Message: "OK"}, nil
}
