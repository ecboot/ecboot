// commission_reverse_test.go 冲销适配器的交付测试（017 评审 I3/N2——原实现零交付测试）:
// 按 after_sale_order.order_item_id 精确冲销（只冲指向项, 不展开整单）; 售后单不存在兜底全单。
package bootstrap

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// seedRevCommissionFixture 造两项已完成订单 + 各一条已结算佣金 + 一张指向 item1 的售后单。
// 返回 item1/item2 id 与清理函数。
func seedRevCommissionFixture(ctx context.Context, t *gtest.T, buyer int64) (item1, item2 int64, cleanup func()) {
	cleanTables := func() {
		for _, no := range []string{"TF-REV-1", "TF-REV-2"} {
			_, _ = g.DB().Exec(ctx, "DELETE c FROM commission_record c JOIN trade_order o ON c.order_no=o.order_no WHERE o.order_no=?", no)
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE order_no=?", no)
			_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_item WHERE order_no=?", no)
			_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", no)
		}
		_, _ = g.DB().Exec(ctx, "DELETE FROM account_log WHERE user_id=?", buyer)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_account WHERE user_id=?", buyer)
	}
	cleanTables()

	var ids []int64
	for i, no := range []string{"TF-REV-1", "TF-REV-2"} {
		res, err := g.DB().Exec(ctx,
			"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
				"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail) "+
				"VALUES(?,?,40,?,0,?,'CNY','冲销测试','13800000000','浙江省','杭州市','T路')",
			no, buyer, "100.00", "100.00")
		t.AssertNil(err)
		oid, _ := res.LastInsertId()
		ires, err := g.DB().Exec(ctx,
			"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,"+
				"quantity,original_price,price,pay_amount) VALUES(?,?,1,?,'TF-REV-SKU','冲销商品','冲销商品','{}',1,'100.00','100.00','100.00')",
			no, oid, buyer+int64(i)+1)
		t.AssertNil(err)
		iid, _ := ires.LastInsertId()
		ids = append(ids, iid)
		// 已结算佣金 10.00（受益人=buyer 自身, 简化断言对象）
		_, err = g.DB().Exec(ctx,
			"INSERT INTO commission_record(order_no,order_item_id,beneficiary_user_id,level,base_amount,rate,amount,status,settle_time) "+
				"VALUES(?,?,?,1,'100.00','10.00','10.00',2,NOW())", no, iid, buyer)
		t.AssertNil(err)
	}
	// 售后单指向第一项
	oid1, err := g.DB().GetValue(ctx, "SELECT id FROM trade_order WHERE order_no='TF-REV-1'")
	t.AssertNil(err)
	_, err = g.DB().Exec(ctx,
		"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,currency,quantity,reason,refund_amount) "+
			"VALUES('TF-REV-AS',?,'TF-REV-1',?,?,2,50,'CNY',1,'测试','10.00')",
		oid1, ids[0], buyer)
	t.AssertNil(err)

	return ids[0], ids[1], cleanTables
}

// TestCommissionReverseAdapterByItem I3 守卫: 售后单指向 item1 → 只冲 item1, item2 佣金不动。
func TestCommissionReverseAdapterByItem(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		buyer := int64(997001)
		item1, item2, cleanup := seedRevCommissionFixture(ctx, t, buyer)
		defer cleanup()
		seed := func() { _, _ = g.DB().Exec(ctx,
			"INSERT INTO user_account(user_id,balance,frozen) VALUES(?,'100.00','0.00') ON DUPLICATE KEY UPDATE balance='100.00', frozen='0.00'", buyer) }
		seed()

		ad := commissionReverseAdapter{}
		ad.ReverseForAfterSale(ctx, "TF-REV-1", "TF-REV-AS", 1000) // refundFen 是参考量, 不影响按项全额

		r1, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE order_item_id=? AND reversal_of_id IS NOT NULL", item1)
		t.AssertNil(err)
		t.Assert(r1.Int(), 1) // item1 被冲销
		r2, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE order_item_id=? AND reversal_of_id IS NOT NULL", item2)
		t.AssertNil(err)
		t.Assert(r2.Int(), 0) // item2 佣金不动（原整单展开实现会误冲）
		acc, _ := g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", buyer)
		t.Assert(acc["balance"].String(), "90.00") // 只扣 item1 的 10.00

		// 兜底: 售后单不存在 → 全单冲销（防御分支, 语义=未知售后范围时保守处理）
		ad.ReverseForAfterSale(ctx, "TF-REV-2", "TF-REV-NOT-EXIST", 0)
		r3, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE order_item_id=? AND reversal_of_id IS NOT NULL", item2)
		t.AssertNil(err)
		t.Assert(r3.Int(), 1) // 兜底全单冲销生效
	})
}
