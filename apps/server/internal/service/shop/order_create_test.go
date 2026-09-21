// order_create_test.go 下单九步事务（012 修复轮补测）。
// 背景: 批次 06 对 `OrderLogicImpl.Create` **零测试覆盖**, 且收口冒烟 6 项恰好绕开下单端点,
// 致两处快照写入缺陷长期存在——① 订单项漏写 4 个非空列（sku_no/spu_name/sku_name/original_price）
// ② 状态流水用不存在的 `operator` 列（与 I6 同根因, 评审仅在 Cancel 中捕获）。
// 两处均使下单事务报错回滚, 即 C 端 `/shop/orders` 实际不可用。
package shop

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// seedUserAddress 建测试收货地址（下单快照来源）。
func seedUserAddress(ctx context.Context, t *gtest.T, userId int64) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO user_address(user_id,receiver_name,receiver_phone,province,city,district,detail_address,is_default) "+
			"VALUES(?,?,?,?,?,?,?,1)",
		userId, "测试收货人", "13800000000", "浙江省", "杭州市", "西湖区", "T路1号")
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// cleanupOrderCreate 清理下单相关数据（流水/订单项/订单/地址）。
func cleanupOrderCreate(ctx context.Context, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE l FROM trade_order_log l JOIN trade_order o ON l.order_id=o.id WHERE o.user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_address WHERE user_id=?", userId)
}

// TestOrderCreateSnapshot 下单成功（FR-001）: 订单头/订单项快照（非空列写全）/状态流水（操作者两列）/
// 库存锁定（只动 locked, 不动 total）/金额账本。
func TestOrderCreateSnapshot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-CRT-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-CRT-TOKEN-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 2,
		})
		t.AssertNil(err)
		t.Assert(out.OrderNo != "", true)
		t.Assert(out.PayAmount, "20.00") // 10.00 × 2; 测试库无在架促销活动

		// 订单头: 状态 10 + 收货快照 + 幂等 token
		rec, err := g.DB().GetOne(ctx,
			"SELECT id, status, total_amount, pay_amount, receiver_name, request_token FROM trade_order WHERE order_no=?",
			out.OrderNo)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 10)
		t.Assert(rec["total_amount"].String(), "20.00")
		t.Assert(rec["pay_amount"].String(), "20.00")
		t.Assert(rec["receiver_name"].String(), "测试收货人")
		t.Assert(rec["request_token"].String(), "T-CRT-TOKEN-1")

		// 订单项快照: 4 个非空列（sku_no/spu_name/sku_name/original_price）必须写全
		it, err := g.DB().GetOne(ctx,
			"SELECT sku_no, spu_name, sku_name, original_price, price, quantity, pay_amount "+
				"FROM trade_order_item WHERE order_id=?", rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(it["sku_no"].String(), "TF-SKU-001")
		t.Assert(it["spu_name"].String(), "TF-商品")
		t.Assert(it["sku_name"].String(), "TF-商品 黑")
		t.Assert(it["original_price"].String(), "10.00")
		t.Assert(it["price"].String(), "10.00")
		t.Assert(it["quantity"].Int(), 2)
		t.Assert(it["pay_amount"].String(), "20.00")

		// 状态流水: 建单 to_status=10; operator_type=2（表注释: 1系统 2用户 3管理员）
		lg, err := g.DB().GetOne(ctx,
			"SELECT to_status, operator_type, operator_id, remark FROM trade_order_log WHERE order_id=?",
			rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(lg["to_status"].Int(), 10)
		t.Assert(lg["operator_type"].Int(), 2)
		t.Assert(lg["operator_id"].String(), "user:"+fmt.Sprint(uid))
		t.Assert(lg["remark"].String(), "订单创建")

		// 库存锁定: locked +2, total 不变
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 100)
		t.Assert(inv["locked"].Int(), 2)
	})
}

// TestOrderCreateInsufficientStock 库存不足（FR-001）: 锁定条件更新未命中 → 40001, 且整单不落库
// （原实现不判 RowsAffected, 条件未命中仍继续下单 → 静默超卖）。
func TestOrderCreateInsufficientStock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-CRT-2"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)
		// 可售仅 1, 下单 2
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=1, locked=0 WHERE sku_id=?", f.SkuId)

		_, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-CRT-TOKEN-2", AddressId: addrId, SkuId: f.SkuId, Quantity: 2,
		})
		t.Assert(errCode(err), errcode.CodeStockInsufficient)

		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM trade_order WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 0)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["locked"].Int(), 0)
	})
}
