// cart_checkout_test.go 购物车修改三态（I9）与结算试算可用券口（I5）回归。
package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

// stubCouponQuery 券查询口桩（I5: 端口由装配层注入, 测试内替换并还原）。
type stubCouponQuery struct{ calls int }

func (s *stubCouponQuery) UsableForOrder(
	_ context.Context, _ int64, _ string,
) ([]model.UsableCouponBrief, error) {
	s.calls++
	return []model.UsableCouponBrief{{UserCouponId: 7, Name: "测试券", Discount: "20.00"}}, nil
}

// TestCartUpdateItemThreeState checked 三态（I9）: 不传 = 不改。
// 原 controller 把 api 的 bool 包成恒非 nil 的指针（`checked := req.Checked; &checked`）,
// 于是"只改数量"的请求会把勾选静默清掉——结算金额随之变化。
func TestCartUpdateItemThreeState(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "CART-I9-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM cart_item WHERE user_id=?", uid) }()

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)

		itemId, err := NewCartLogic().AddItem(ctx, uid, f.SkuId, 1)
		t.AssertNil(err)
		checked := func() int {
			v, _ := g.DB().GetValue(ctx, "SELECT checked FROM cart_item WHERE id=?", itemId)
			return v.Int()
		}
		t.Assert(checked(), 1) // 加购默认勾选

		// 只改数量（checked 传 nil = 不改）→ 勾选必须保持
		t.AssertNil(NewCartLogic().UpdateItem(ctx, uid, itemId, 3, nil))
		qty, err := g.DB().GetValue(ctx, "SELECT quantity FROM cart_item WHERE id=?", itemId)
		t.AssertNil(err)
		t.Assert(qty.Int(), 3)
		t.Assert(checked(), 1)

		// 显式取消勾选
		no := false
		t.AssertNil(NewCartLogic().UpdateItem(ctx, uid, itemId, 0, &no))
		t.Assert(checked(), 0)

		// 显式勾选
		yes := true
		t.AssertNil(NewCartLogic().UpdateItem(ctx, uid, itemId, 0, &yes))
		t.Assert(checked(), 1)
	})
}

// TestCartCheckoutUsableCoupons 结算试算可用券（I5）: 端口已注入则出参带券; 未注入/查询失败降级为空。
func TestCartCheckoutUsableCoupons(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		defer func(orig ICouponQuery) { CouponQuery = orig }(CouponQuery)

		const h1 = "CART-I5-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM cart_item WHERE user_id=?", uid) }()

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		_, err := NewCartLogic().AddItem(ctx, uid, f.SkuId, 2)
		t.AssertNil(err)

		// 已注入: 出参带可用券, 且按商品总额（元）查询
		stub := &stubCouponQuery{}
		CouponQuery = stub
		res, err := NewCartLogic().Checkout(ctx, uid, model.CheckoutQuery{})
		t.AssertNil(err)
		t.Assert(stub.calls, 1)
		t.Assert(len(res.UsableCoupons), 1)
		t.Assert(res.UsableCoupons[0].UserCouponId, int64(7))
		t.Assert(res.UsableCoupons[0].Discount, "20.00")
		t.Assert(res.Amount.TotalAmount, "20.00")

		// 未注入（装配缺失）: 降级为空且不报错
		CouponQuery = nil
		res, err = NewCartLogic().Checkout(ctx, uid, model.CheckoutQuery{})
		t.AssertNil(err)
		t.Assert(len(res.UsableCoupons), 0)
	})
}
