package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
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

// seedLogistics 建物流公司（用后清理）。
func seedLogisticsForOrder(ctx context.Context, t *gtest.T, code string, status int) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM logistics_company WHERE code=?", code)
	_, err := g.DB().Exec(ctx,
		"INSERT INTO logistics_company(code,name,tracking_rule,sort,status) VALUES(?,?,'',0,?)",
		code, code+"-name", status)
	t.AssertNil(err)
}

// TestAdminOrderDeliver 发货（FR-010）: 20→30 + 物流落库; 非 20 → 40006; 停用物流公司 → 拒绝。
func TestAdminOrderDeliver(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ADM-DLV-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		const lcOn, lcOff = "T-DLV-ON", "T-DLV-OFF"
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM logistics_company WHERE code IN (?,?)", lcOn, lcOff)
		}()
		seedLogisticsForOrder(ctx, t, lcOn, 1)
		seedLogisticsForOrder(ctx, t, lcOff, 0)

		// 待发货（20）→ 30
		id := seedOrder(ctx, t, uid, "T-DLV-OK", 20)
		defer cleanupOrder(ctx, t, "T-DLV-OK")
		t.AssertNil(NewOrderLogic().Deliver(ctx, "T-DLV-OK", lcOn, "SF123456", "admin:1"))
		rec, err := g.DB().GetOne(ctx,
			"SELECT status, deliver_company, deliver_no FROM trade_order WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 30)
		t.Assert(rec["deliver_company"].String(), lcOn)
		t.Assert(rec["deliver_no"].String(), "SF123456")

		// 非待发货（10）→ 40006
		seedOrder(ctx, t, uid, "T-DLV-BAD", 10)
		defer cleanupOrder(ctx, t, "T-DLV-BAD")
		err = NewOrderLogic().Deliver(ctx, "T-DLV-BAD", lcOn, "SF1", "admin:1")
		t.Assert(errCode(err), errcode.CodeStatusNotAllowed)

		// 停用物流公司 → 10001（拒绝）
		seedOrder(ctx, t, uid, "T-DLV-OFF", 20)
		defer cleanupOrder(ctx, t, "T-DLV-OFF")
		err = NewOrderLogic().Deliver(ctx, "T-DLV-OFF", lcOff, "SF1", "admin:1")
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}

// TestAdminOrderListAndRemark 后台列表与备注（FR-009/012）。
func TestAdminOrderListAndRemark(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ADM-LST-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		seedOrder(ctx, t, uid, "T-ADM-L1", 20)
		seedOrder(ctx, t, uid, "T-ADM-L2", 30)
		defer cleanupOrder(ctx, t, "T-ADM-L1")
		defer cleanupOrder(ctx, t, "T-ADM-L2")

		// 状态筛选（20）
		res, err := NewOrderLogic().AdminList(ctx, model.AdminOrderQuery{
			Status: 20, OrderNo: "T-ADM-L1", PageReq: model.PageReq{Page: 1, PageSize: 10},
		})
		t.AssertNil(err)
		t.Assert(res.Total, 1)
		t.Assert(res.List[0].OrderNo, "T-ADM-L1")

		// 备注落库
		t.AssertNil(NewOrderLogic().SellerRemark(ctx, "T-ADM-L1", "内部备注x", "admin:1"))
		v, err := g.DB().GetValue(ctx, "SELECT seller_remark FROM trade_order WHERE order_no='T-ADM-L1'")
		t.AssertNil(err)
		t.Assert(v.String(), "内部备注x")
		// 后台详情可见该备注（C 端 controller 不映射该字段——买家不可见）
		d, err := NewOrderLogic().AdminDetail(ctx, "T-ADM-L1")
		t.AssertNil(err)
		t.Assert(d.SellerRemark, "内部备注x")
	})
}
