package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/shop"
)

// PayCreate 发起支付（返回渠道唤起参数）
func (c *ControllerV1) PayCreate(ctx context.Context, req *v1.PayCreateReq) (res *v1.PayCreateRes, err error) {
	out, err := shop.NewPayLogic().Create(ctx, middleware.CtxUserIdFrom(ctx), req.OrderNo, req.PayChannel)
	if err != nil {
		return nil, err
	}
	return &v1.PayCreateRes{PayNo: out.PayNo, ChannelParams: out.ChannelParams}, nil
}
