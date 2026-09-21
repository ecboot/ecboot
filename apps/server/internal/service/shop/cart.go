// cart.go 购物车域——表: cart_item（V4）。
// 规则: 实时联查 SKU 价态（不快照）; 重复加购=数量累加（UNIQUE 兜底）;
// 结算页实时试算: 商品价态 + 优惠试算(券/满减/积分/余额) + 运费; 下架/禁售置灰不阻断。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// ICartLogic 购物车。
type ICartLogic interface {
	Detail(ctx context.Context, userId int64) (*model.CartView, error)
	AddItem(ctx context.Context, userId, skuId int64, quantity int) (int64, error)
	UpdateItem(ctx context.Context, userId, itemId int64, quantity int, checked *bool) error
	RemoveItem(ctx context.Context, userId, itemId int64) error
	// Checkout 结算试算（不改状态; 优惠按 004 契约口径: 先满减后用券, 积分/余额同层分摊）。
	Checkout(ctx context.Context, userId int64, in model.CheckoutQuery) (*model.CheckoutResult, error)
}

// model.AmountBook 金额账本（四优惠构成 + 运费 + 现金应付; 恒等式见 contracts §4）。
