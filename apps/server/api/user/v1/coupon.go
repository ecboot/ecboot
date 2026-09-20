package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	UsableCouponTemplate struct {
		CouponId    string `json:"couponId"`
		Name        string `json:"name" dc:"券名称"`
		Type        int    `json:"type" dc:"1满减 2无门槛"`
		Threshold   string `json:"threshold" dc:"门槛金额"`
		Discount    string `json:"discount" dc:"抵扣金额"`
		ValidDesc   string `json:"validDesc" dc:"有效期描述"`
		CanReceive  bool   `json:"canReceive" dc:"是否可领(限领/存量)"`
	}

	// 可领模板列表
	CouponAvailableListReq struct {
		g.Meta `path:"/coupons/available" method:"GET" summary:"可领优惠券列表"`
		PageReq
	}
	CouponAvailableListRes struct {
		PageRes
		List []UsableCouponTemplate `json:"list"`
	}

	// 领券（防超发/限领, 幂等语义由服务端保证）
	CouponReceiveReq struct {
		g.Meta   `path:"/coupons/{couponId}/receive" method:"POST" summary:"领取优惠券"`
		CouponId string `json:"couponId" v:"required" dc:"券模板ID"`
	}
	CouponReceiveRes struct {
		UserCouponId string `json:"userCouponId" dc:"用户券ID"`
	}

	// 我的券（按状态）
	MyCouponItem struct {
		UserCouponId string `json:"userCouponId"`
		Name         string `json:"name"`
		Threshold    string `json:"threshold" dc:"门槛"`
		Discount     string `json:"discount" dc:"抵扣"`
		ExpireTime   string `json:"expireTime" dc:"过期时间"`
		Status       int    `json:"status" dc:"1未使用 2已使用 3已过期 4已退回"`
	}
	MyCouponListReq struct {
		g.Meta  `path:"/coupons" method:"GET" summary:"我的优惠券"`
		Status  int `json:"status" dc:"状态筛选" d:"1"`
		PageReq
	}
	MyCouponListRes struct {
		PageRes
		List []MyCouponItem `json:"list"`
	}
)
