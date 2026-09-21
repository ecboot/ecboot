package shop

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// ---- 支付 fixture（012） ----

// seedPayFixture 建待付款订单 + 订单项 + 库存行（total=10, locked=2）, 返回 (orderId, skuId)。
func seedPayFixture(ctx context.Context, t *gtest.T, userId int64, orderNo string, payFen int64) (int64, int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_item WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM pay_order WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM pay_callback_log WHERE pay_no LIKE 'PAY-T-%'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)

	// 库存行（sku_id 用哨兵大数避免与真实 SKU 冲突）
	res, err := g.DB().Exec(ctx,
		"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,10,2,0) ON DUPLICATE KEY UPDATE total=10,locked=2",
		"990000001")
	t.AssertNil(err)
	_ = res

	amount := float64(payFen) / 100.0
	res, err = g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail) "+
			"VALUES(?,?,10,?,0,?,'CNY','测试','13800000000','浙江省','杭州市','T路1号')",
		orderNo, userId, amount, amount)
	t.AssertNil(err)
	orderId, _ := res.LastInsertId()

	_, err = g.DB().Exec(ctx,
		"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,quantity,original_price,price,promotion_amount,pay_amount) "+
			"VALUES(?,?,1,990000001,'T-SKU','测试商品','','{}',2,?,?,0,?)",
		orderNo, orderId, amount, amount, amount)
	t.AssertNil(err)
	return orderId, 990000001
}

func cleanupPayFixture(ctx context.Context, t *gtest.T, orderNo string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_log WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_item WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM pay_order WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
}

// TestPayCreate 发起支付（FR-005）: 待付款订单成功（返回 mock 唤起参数）; 非待付款/他人 → 拒绝。
func TestPayCreate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-C-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		_, _ = seedPayFixture(ctx, t, uid, "T-PAY-OK", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-OK")

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-OK", 1)
		t.AssertNil(err)
		t.Assert(out.PayNo != "", true)
		t.Assert(len(out.ChannelParams) > 0, true) // mock 唤起参数

		// 非待付款（已支付 20）→ 40006
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET status=20 WHERE order_no='T-PAY-OK'")
		_, err = NewPayLogic().Create(ctx, uid, "T-PAY-OK", 1)
		t.Assert(errCode(err), errcode.CodeStatusNotAllowed)
	})
}

// TestPayNotifyIdempotent 支付回调（FR-006/007）: 成功推进（订单 20 + 库存核销）; 重复回调不二次推进。
func TestPayNotifyIdempotent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-N-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		orderId, skuId := seedPayFixture(ctx, t, uid, "T-PAY-N1", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-N1")
		// 库存基线: total=10, locked=2
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=2 WHERE sku_id=?", skuId)

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-N1", 1)
		t.AssertNil(err)

		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-1", "amountFen": 10000, "success": true,
		})
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body))

		// 订单 → 20（待发货）
		st, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(st.Int(), 20)
		// 库存核销: total 10→8（扣 2）, locked 2→0
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 8)
		t.Assert(inv["locked"].Int(), 0)
		// 支付单 → 20
		pst, err := g.DB().GetValue(ctx, "SELECT status FROM pay_order WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.Assert(pst.Int(), 20)
		// 回调留档
		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM pay_callback_log WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.AssertGE(n.Int(), 1)

		// 重复回调: 幂等（库存不再扣减, 订单状态不变）
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body))
		inv, err = g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 8) // 未二次核销
	})
}

// TestPayNotifyAmountMismatch 金额不符（FR-007）: 拒绝且不推进（留档）。
func TestPayNotifyAmountMismatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-M-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		orderId, _ := seedPayFixture(ctx, t, uid, "T-PAY-M1", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-M1")

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-M1", 1)
		t.AssertNil(err)
		// 回调金额 1 分（应付 10000 分）→ 拒绝
		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-M", "amountFen": 1, "success": true,
		})
		err = NewPayLogic().HandlePayNotify(ctx, "mock", body)
		t.AssertNE(err, nil)
		// 订单未推进
		st, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(st.Int(), 10)
		// 留档（成败均留）
		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM pay_callback_log WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.AssertGE(n.Int(), 1)
	})
}
