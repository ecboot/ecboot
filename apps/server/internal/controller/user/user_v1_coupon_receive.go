package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// CouponReceive 领取优惠券（同事务防超发+限领）
func (c *ControllerV1) CouponReceive(ctx context.Context, req *v1.CouponReceiveReq) (res *v1.CouponReceiveRes, err error) {
	cid, err := parseID(req.CouponId)
	if err != nil {
		return nil, err
	}
	id, err := user.Receive(ctx, middleware.CtxUserIdFrom(ctx), cid)
	if err != nil {
		return nil, err
	}
	return &v1.CouponReceiveRes{UserCouponId: fmtID(id)}, nil
}
