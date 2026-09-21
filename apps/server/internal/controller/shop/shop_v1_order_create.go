package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// OrderCreate 创建订单（幂等 requestToken; 玩法上下文互斥）
func (c *ControllerV1) OrderCreate(ctx context.Context, req *v1.OrderCreateReq) (res *v1.OrderCreateRes, err error) {
	addrId, err := parseID(req.AddressId)
	if err != nil {
		return nil, err
	}
	var couponId, skuId, teamId, flashId, bargainId int64
	if couponId, err = optID(req.UserCouponId); err != nil {
		return nil, err
	}
	if skuId, err = optID(req.SkuId); err != nil {
		return nil, err
	}
	if teamId, err = optID(req.GroupBuyTeamId); err != nil {
		return nil, err
	}
	if flashId, err = optID(req.FlashSaleItemId); err != nil {
		return nil, err
	}
	if bargainId, err = optID(req.BargainRecordId); err != nil {
		return nil, err
	}
	cartIds := make([]int64, 0, len(req.CartItemIds))
	for _, s := range req.CartItemIds {
		id, perr := parseID(s)
		if perr != nil {
			return nil, perr
		}
		cartIds = append(cartIds, id)
	}
	out, err := shop.NewOrderLogic().Create(ctx, middleware.CtxUserIdFrom(ctx), model.OrderCreateInput{
		RequestToken: req.RequestToken, AddressId: addrId, UserCouponId: couponId,
		UsePoint: req.UsePoint, UseAccount: req.UseAccount, UserRemark: req.UserRemark,
		CartItemIds: cartIds, GroupBuyTeamId: teamId, SkuId: skuId, Quantity: req.Quantity,
		FlashSaleItemId: flashId, BargainRecordId: bargainId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.OrderCreateRes{OrderNo: out.OrderNo, PayAmount: out.PayAmount}, nil
}

// optID 可选 ID（空串=0; shop 渠道共享）。
func optID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return parseID(s)
}
