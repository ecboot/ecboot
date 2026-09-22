// order_confirm_hook_test.go 确认收货 → 佣金计提事件的投递守卫（017 评审 I6）:
// Confirm 在 30→40 迁移成功后必须投递 ICommissionSettle（spy 捕获），未装配时走告警降级。
package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

type settleSpy struct {
	calls    int
	lastNo   string
	installed bool
}

func (s *settleSpy) OnOrderConfirmed(_ context.Context, orderNo string) {
	s.calls++
	s.lastNo = orderNo
}

// TestConfirmDispatchesCommissionSettle 确认收货投递计提事件（I6）。
func TestConfirmDispatchesCommissionSettle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "CFM-HOOK-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		// 造一笔 30 态订单（下单后手工推状态到 30）
		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-CFM-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 1,
		})
		t.AssertNil(err)
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET status=30 WHERE order_no=?", out.OrderNo)
		t.Cleanup(func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_log WHERE order_no=?", out.OrderNo)
		})

		// 装配 spy
		old := CommissionSettle
		spy := &settleSpy{installed: true}
		CommissionSettle = spy
		defer func() { CommissionSettle = old }()

		t.AssertNil(NewOrderLogic().Confirm(ctx, uid, out.OrderNo))
		t.Assert(spy.calls, 1)
		t.Assert(spy.lastNo, out.OrderNo)
	})
}
