package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	CartItem struct {
		ItemId    string            `json:"itemId" dc:"购物车项ID"`
		SkuId     string            `json:"skuId"`
		SpuName   string            `json:"spuName" dc:"商品名"`
		Specs     map[string]string `json:"specs" dc:"规格"`
		Image     string            `json:"image" dc:"SKU图"`
		Price     string            `json:"price" dc:"现价(元)"`
		Sellable  bool              `json:"sellable" dc:"是否可售"`
		Quantity  int               `json:"quantity" dc:"数量"`
		Checked   bool              `json:"checked" dc:"结算勾选"`
	}

	CartDetailReq struct {
		g.Meta `path:"/cart" method:"GET" summary:"购物车详情"`
	}
	CartDetailRes struct {
		List []CartItem `json:"list"`
	}

	CartAddItemReq struct {
		g.Meta `path:"/cart/items" method:"POST" summary:"加入购物车"`
		SkuId  string `json:"skuId" v:"required" dc:"SKU ID"`
		Quantity int  `json:"quantity" v:"required|min:1" dc:"数量" d:"1"`
	}
	CartAddItemRes struct {
		ItemId string `json:"itemId" dc:"购物车项ID"`
	}

	CartUpdateItemReq struct {
		g.Meta   `path:"/cart/items/{itemId}" method:"PUT" summary:"修改购物车项"`
		ItemId   string `json:"itemId" v:"required" dc:"购物车项ID"`
		Quantity int    `json:"quantity" dc:"数量"`
		Checked  bool   `json:"checked" dc:"勾选"`
	}
	CartUpdateItemRes struct {
		Success bool `json:"success"`
	}

	CartRemoveItemReq struct {
		g.Meta `path:"/cart/items/{itemId}" method:"DELETE" summary:"移除购物车项"`
		ItemId string `json:"itemId" v:"required" dc:"购物车项ID"`
	}
	CartRemoveItemRes struct {
		Success bool `json:"success"`
	}

	// 结算试算（勾选项: 商品价态+优惠明细+运费+应付）
	CartCheckoutReq struct {
		g.Meta     `path:"/cart/checkout" method:"GET" summary:"结算试算"`
		AddressId  string `json:"addressId" dc:"地址ID(算运费)"`
		CouponId   string `json:"couponId" dc:"使用的用户券"`
		UsePoint   bool   `json:"usePoint" dc:"是否用积分抵扣"`
		UseAccount bool   `json:"useAccount" dc:"是否用佣金余额抵扣"`
	}
	CheckoutAmount struct {
		TotalAmount        string `json:"totalAmount" dc:"商品总额"`
		CouponAmount       string `json:"couponAmount" dc:"券抵扣"`
		FullReductionAmount string `json:"fullReductionAmount" dc:"满减"`
		PointAmount        string `json:"pointAmount" dc:"积分抵扣"`
		AccountAmount      string `json:"accountAmount" dc:"余额抵扣"`
		FreightAmount      string `json:"freightAmount" dc:"运费"`
		PayAmount          string `json:"payAmount" dc:"应付(现金部分)"`
	}
	CartCheckoutRes struct {
		Items       []CartItem       `json:"items" dc:"结算商品"`
		Amount      CheckoutAmount   `json:"amount"`
		UsableCoupons []UsableCoupon `json:"usableCoupons" dc:"可用券列表"`
		Errors      []string         `json:"errors" dc:"失效商品提示"`
	}
	UsableCoupon struct {
		UserCouponId string `json:"userCouponId"`
		Name         string `json:"name"`
		Threshold    string `json:"threshold" dc:"门槛"`
		Discount     string `json:"discount" dc:"抵扣"`
	}
)
