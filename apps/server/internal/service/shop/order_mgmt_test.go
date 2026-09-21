package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// ---- 订单管理面 fixture（012） ----

// seedOrder 建测试订单（指定状态）, 返回 id。
func seedOrder(ctx context.Context, t *gtest.T, userId int64, orderNo string, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail,deliver_company,deliver_no,seller_remark) "+
			"VALUES(?,?,?,100.00,100.00,'CNY','测试','13800000000','浙江省','杭州市','T路1号','','','')",
		orderNo, userId, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupOrder(ctx context.Context, t *gtest.T, orderNo string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_log WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
}

// TestOrderConfirm C 端确认收货（FR-004）: 30→40 + 日志; 非 30 → 40006; 他人 → 40005。
func TestOrderConfirm(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			h1 = "ORD-CFM-1"
			h2 = "ORD-CFM-2"
		)
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		uid := seedPointUser(ctx, t, h1)
		other := seedPointUser(ctx, t, h2)

		// 待收货（30）→ 确认收货 → 已完成（40）
		id := seedOrder(ctx, t, uid, "T-CFM-OK", 30)
		defer cleanupOrder(ctx, t, "T-CFM-OK")
		t.AssertNil(NewOrderLogic().Confirm(ctx, uid, "T-CFM-OK"))
		st, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(st.Int(), 40)
		// 状态流水
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM trade_order_log WHERE order_no=? AND to_status=40", "T-CFM-OK")
		t.AssertNil(err)
		t.Assert(n.Int(), 1)

		// 非待收货（10）→ 40006
		seedOrder(ctx, t, uid, "T-CFM-BAD", 10)
		defer cleanupOrder(ctx, t, "T-CFM-BAD")
		err = NewOrderLogic().Confirm(ctx, uid, "T-CFM-BAD")
		t.Assert(errCode(err), errcode.CodeStatusNotAllowed)

		// 他人订单 → 40005（不可见语义）
		err = NewOrderLogic().Confirm(ctx, other, "T-CFM-OK")
		t.Assert(errCode(err), errcode.CodeOrderNotFound)
	})
}
