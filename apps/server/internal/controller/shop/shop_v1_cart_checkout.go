package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// CartCheckout 结算试算（金额账本勾稽; 积分抵扣由 012 前置修复承接）
func (c *ControllerV1) CartCheckout(ctx context.Context, req *v1.CartCheckoutReq) (res *v1.CartCheckoutRes, err error) {
	var addrId, couponId int64
	if req.AddressId != "" {
		if addrId, err = parseID(req.AddressId); err != nil {
			return nil, err
		}
	}
	if req.CouponId != "" {
		if couponId, err = parseID(req.CouponId); err != nil {
			return nil, err
		}
	}
	out, err := shop.NewCartLogic().Checkout(ctx, middleware.CtxUserIdFrom(ctx), model.CheckoutQuery{
		AddressId: addrId, CouponId: couponId, UsePoint: req.UsePoint, UseAccount: req.UseAccount,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.CartCheckoutRes{
		Errors: out.Errors,
		Amount: v1.CheckoutAmount{
			TotalAmount: out.Amount.TotalAmount, CouponAmount: out.Amount.CouponAmount,
			FullReductionAmount: out.Amount.FullReductionAmount, PointAmount: out.Amount.PointAmount,
			AccountAmount: out.Amount.AccountAmount,
			FreightAmount: out.Amount.FreightAmount, PayAmount: out.Amount.PayAmount,
		},
		Items:         make([]v1.CartItem, 0, len(out.Items)),
		UsableCoupons: make([]v1.UsableCoupon, 0, len(out.UsableCoupons)),
	}
	for _, it := range out.Items {
		res.Items = append(res.Items, v1.CartItem{
			ItemId: fmtID(it.ItemId), SkuId: fmtID(it.SkuId), SpuName: it.SpuName,
			Specs: it.Specs, Image: it.Image, Price: it.Price,
			Sellable: it.Sellable, Quantity: it.Quantity, Checked: it.Checked,
		})
	}
	for _, cp := range out.UsableCoupons {
		res.UsableCoupons = append(res.UsableCoupons, v1.UsableCoupon{
			UserCouponId: fmtID(cp.UserCouponId), Name: cp.Name, Discount: cp.Discount,
		})
	}
	return res, nil
}
