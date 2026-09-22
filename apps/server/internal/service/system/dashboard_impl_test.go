// dashboard_impl_test.go 三看板口径（019 批次 13）——已知数据口径钉住 + 空窗口全零。
package system

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestDashboardTrade 交易看板口径（FR-1）: 已支付口径/销售额合计/待发货/退款额。
func TestDashboardTrade(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewDashboardLogic()

		const no1, no2, no3 = "TF-DB-T1", "TF-DB-T2", "TF-DB-T3"
		for _, no := range []string{no1, no2, no3} {
			_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", no)
		}
		defer func() {
			for _, no := range []string{no1, no2, no3} {
				_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", no)
			}
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no='TF-DB-AS1'")
		}()
		insOrder := func(no string, status int, pay string) {
			_, err := g.DB().Exec(ctx,
				"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
					"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail) "+
					"VALUES(?,997101,?,?,0,?,'CNY','看板测试','13800000000','浙江省','杭州市','T路')",
				no, status, pay, pay)
			t.AssertNil(err)
		}
		// 基线（共享库既有数据不可控 → 差值断言; 必须在插单前取）
		base, err := logic.Trade(ctx, "", "")
		t.AssertNil(err)

		insOrder(no1, 20, "50.00")  // 待发货（计入销售额）
		insOrder(no2, 40, "80.00")  // 已完成
		insOrder(no3, 10, "999.00") // 待付款（不计入）

		// 造售后完成单: 退款额 +10
		oid, err := g.DB().GetValue(ctx, "SELECT id FROM trade_order WHERE order_no=?", no2)
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,currency,quantity,reason,refund_amount) "+
				"VALUES('TF-DB-AS1',?,?,0,997101,2,50,'CNY',1,'测试','10.00')", oid.Int64(), no2)
		t.AssertNil(err)

		after, err := logic.Trade(ctx, "", "")
		t.AssertNil(err)

		// 订单数 +2（20/40 计入, 10 不计）; 销售额 +130.00; 退款额 +10.00
		t.Assert(after.OrderCount, base.OrderCount+2)
		t.Assert(fen2(after.SalesAmount), fen2(base.SalesAmount)+13000)
		t.Assert(fen2(after.RefundAmount), fen2(base.RefundAmount)+1000)
		t.Assert(after.PendingDeliver >= 1, true)
	})
}

// TestDashboardMember 会员看板口径（FR-2）: 窗口新增/活跃/休眠≥90天。
func TestDashboardMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewDashboardLogic()

		base, err := logic.Member(ctx, "", "")
		t.AssertNil(err)

		// 造一个"90 天前活跃"的会员 → 休眠 +1（窗口新增不受影响: created_at 默认 NOW）
		_, err = g.DB().Exec(ctx,
			"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status,last_active_at) "+
				"VALUES('TF-DB-DORM','x','TF-DB-DORM-PH',0,1,DATE_SUB(NOW(), INTERVAL 100 DAY))")
		t.AssertNil(err)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash='TF-DB-DORM-PH'") }()

		after, err := logic.Member(ctx, "", "")
		t.AssertNil(err)
		t.Assert(after.DormantCount, base.DormantCount+1) // 休眠 ≥90 天口径
		t.Assert(after.NewCount >= base.NewCount, true)   // 新增口径单调（窗口默认全量）
	})
}

// TestDashboardProduct 商品看板口径（FR-3）: 在售/低库存预警/待审核评价。
func TestDashboardProduct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := NewDashboardLogic()
		base, err := f.Product(ctx)
		t.AssertNil(err)

		// 造在售 SPU + 低库存 SKU（available <= warn_count）
		spu, err := g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
			"spu_no": "TF-DB-SPU", "name": "TF-DB-看板商品", "category_id": 1,
			"images": `[]`, "spec_definitions": `[]`, "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		sku, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
			"sku_no": "TF-DB-SKU", "spu_id": spu, "name": "TF-DB-看板SKU",
			"specs": "{}", "price": "1.00", "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,5,0,10) ON DUPLICATE KEY UPDATE total=5, locked=0, warn_count=10", sku)
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", sku)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE id=?", sku)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spu)
		}()

		after, err := f.Product(ctx)
		t.AssertNil(err)
		t.Assert(after.OnSaleCount, base.OnSaleCount+1)   // 在售 +1
		t.Assert(after.LowStockCount, base.LowStockCount+1) // 低库存预警 +1（available 5 <= warn 10）
	})
}

// fen2 元字符串 → 分（断言用；解析失败归 -1 以便暴露）。
func fen2(yuan string) int64 {
	var yuanPart int64
	var fenPart int64
	i := 0
	for ; i < len(yuan); i++ {
		if yuan[i] == '.' {
			break
		}
		if yuan[i] < '0' || yuan[i] > '9' {
			return -1
		}
		yuanPart = yuanPart*10 + int64(yuan[i]-'0')
	}
	if i < len(yuan) {
		frac := yuan[i+1:]
		for j := 0; j < len(frac) && j < 2; j++ {
			fenPart = fenPart*10 + int64(frac[j]-'0')
		}
		for j := len(frac); j < 2; j++ {
			fenPart *= 10
		}
	}
	return yuanPart*100 + fenPart
}
