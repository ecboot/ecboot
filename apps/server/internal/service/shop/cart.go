// cart.go 购物车域——表: cart_item（V4）。
// 规则: 实时联查 SKU 价态（不快照）; 重复加购=数量累加（UNIQUE 兜底）;
// 结算页实时试算: 商品价态 + 优惠试算(券/满减/积分/余额) + 运费; 下架/禁售置灰不阻断。
package shop

import "context"

// ICartLogic 购物车。
type ICartLogic interface {
	Detail(ctx context.Context, userId int64) (*CartView, error)
	AddItem(ctx context.Context, userId, skuId int64, quantity int) (int64, error)
	UpdateItem(ctx context.Context, userId, itemId int64, quantity int, checked *bool) error
	RemoveItem(ctx context.Context, userId, itemId int64) error
	// Checkout 结算试算（不改状态; 优惠按 004 契约口径: 先满减后用券, 积分/余额同层分摊）。
	Checkout(ctx context.Context, userId int64, in CheckoutQuery) (*CheckoutResult, error)
}

type CartView struct {
	Items []CartLine `json:"items"`
}

type CartLine struct {
	ItemId   int64             `json:"itemId"`
	SkuId    int64             `json:"skuId"`
	SpuName  string            `json:"spuName"`
	Specs    map[string]string `json:"specs"`
	Image    string            `json:"image"`
	Price    string            `json:"price"`
	Sellable bool              `json:"sellable"`
	Quantity int               `json:"quantity"`
	Checked  bool              `json:"checked"`
}

type CheckoutQuery struct {
	AddressId  int64
	CouponId   int64
	UsePoint   bool
	UseAccount bool
}

type CheckoutResult struct {
	Items         []CartLine          `json:"items"`
	Amount        AmountBook          `json:"amount"`
	UsableCoupons []UsableCouponBrief `json:"usableCoupons"`
	Errors        []string            `json:"errors" dc:"失效商品提示"`
}

// AmountBook 金额账本（四优惠构成 + 运费 + 现金应付; 恒等式见 contracts §4）。
type AmountBook struct {
	TotalAmount         string `json:"totalAmount"`
	CouponAmount        string `json:"couponAmount"`
	FullReductionAmount string `json:"fullReductionAmount"`
	PointAmount         string `json:"pointAmount"`
	PointUsed           int    `json:"pointUsed"`
	AccountAmount       string `json:"accountAmount" dc:"余额抵扣(现金应付=payAmount-accountAmount)"`
	FreightAmount       string `json:"freightAmount"`
	PayAmount           string `json:"payAmount" dc:"实付(含余额抵扣前口径=total-promotion+freight)"`
}

type UsableCouponBrief struct {
	UserCouponId int64  `json:"userCouponId"`
	Name         string `json:"name"`
	Discount     string `json:"discount"`
}
