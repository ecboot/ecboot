package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// PayStatus 支付状态查询
func (c *ControllerV1) PayStatus(ctx context.Context, req *v1.PayStatusReq) (res *v1.PayStatusRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewPayLogic().Status(ctx, userId, req.PayNo)
	if err != nil {
		return nil, err
	}
	return &v1.PayStatusRes{Status: out.Status, PaidAt: out.PaidAt}, nil
}
