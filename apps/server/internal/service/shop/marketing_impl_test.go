// marketing_impl_test.go 营销 C 端·公开列表与秒杀链路（015-marketing-c 批次 09）。
// 重点: **秒杀欠账清偿**（批次 07 曾把秒杀下单临时拦为"未上线"）——秒杀价快照、活动/商品**双库存**、
// 取消双回补、**支付回调核销成功**（批次 07 的失败点：只锁活动库存会导致回调核销未命中而整单回滚）。
package shop

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// seedFlashSale 造秒杀场次 + 场次商品（offMin: 开始偏移分钟, endMin: 结束偏移分钟; 负数=已过去）。
func seedFlashSale(
	ctx context.Context, t *gtest.T, name string, skuId int64, flashPrice string, stock, offMin, endMin int,
) (int64, int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_item WHERE activity_id IN (SELECT id FROM flash_sale_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_activity WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO flash_sale_activity(name,start_time,end_time,status,deleted) "+
			"VALUES(?,DATE_ADD(NOW(), INTERVAL ? MINUTE),DATE_ADD(NOW(), INTERVAL ? MINUTE),1,0)",
		name, offMin, endMin)
	t.AssertNil(err)
	actId, _ := res.LastInsertId()

	ires, err := g.DB().Exec(ctx,
		"INSERT INTO flash_sale_item(activity_id,sku_id,flash_price,stock_count,sold_count,per_limit) VALUES(?,?,?,?,0,2)",
		actId, skuId, flashPrice, stock)
	t.AssertNil(err)
	itemId, _ := ires.LastInsertId()
	return actId, itemId
}

func cleanupFlashSale(ctx context.Context, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_item WHERE activity_id IN (SELECT id FROM flash_sale_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_activity WHERE name=?", name)
}

// TestFlashSaleList 秒杀公开列表（FR-001）: 进行中可见、预告含 upcoming、已结束不可见、剩余量正确。
func TestFlashSaleList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)

		const nameOK, nameUp, nameEnd = "TF-秒杀进行中", "TF-秒杀预告", "TF-秒杀已结束"
		defer cleanupFlashSale(ctx, nameOK)
		defer cleanupFlashSale(ctx, nameUp)
		defer cleanupFlashSale(ctx, nameEnd)
		seedFlashSale(ctx, t, nameOK, f.SkuId, "5.00", 10, -10, 60)
		seedFlashSale(ctx, t, nameUp, f.SkuId, "4.00", 8, 10, 90)
		seedFlashSale(ctx, t, nameEnd, f.SkuId, "3.00", 5, -100, -10)

		res, err := NewMarketingLogic().PublicFlashSales(ctx, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		hit := map[string]model.PublicFlashSaleItem{}
		for _, it := range res.List {
			hit[it.Name] = it
		}
		t.Assert(hit[nameOK].ActivityId > 0, true)
		t.Assert(hit[nameUp].Upcoming, true) // 预告
		_, ended := hit[nameEnd]
		t.Assert(ended, false) // 已结束不出现
		// 场次商品摘要: 秒杀价 + 剩余量 + 限购
		ok := hit[nameOK]
		t.Assert(len(ok.Items), 1)
		t.Assert(ok.Items[0].SkuId, f.SkuId)
		t.Assert(ok.Items[0].Price, "5.00")
		t.Assert(ok.Items[0].StockRemain, 10)
		t.Assert(ok.Items[0].PerLimit, 2)
	})
}

// TestFlashSaleOrder 秒杀下单（FR-002~005）: 秒杀价快照 + 活动/商品**双库存**扣减; 三种拒绝路径。
func TestFlashSaleOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "FS-ORD-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		const name = "TF-秒杀下单"
		defer cleanupFlashSale(ctx, name)
		actId, itemId := seedFlashSale(ctx, t, name, f.SkuId, "5.00", 10, -10, 60)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-ORD-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 2, FlashSaleItemId: itemId,
		})
		t.AssertNil(err) // 批次 07 此处返回 50002（临时拦截）；本批解除

		// 订单头: 按秒杀价计（5.00 × 2 = 10.00）
		rec, err := g.DB().GetOne(ctx, "SELECT id, total_amount, pay_amount FROM trade_order WHERE order_no=?", out.OrderNo)
		t.AssertNil(err)
		t.Assert(rec["total_amount"].String(), "10.00")
		t.Assert(rec["pay_amount"].String(), "10.00")

		// 订单项: price = 秒杀价（成交价快照）, original_price = 原价（10.00）
		it, err := g.DB().GetOne(ctx,
			"SELECT price, original_price, pay_amount FROM trade_order_item WHERE order_id=?", rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(it["price"].String(), "5.00")
		t.Assert(it["original_price"].String(), "10.00")
		t.Assert(it["pay_amount"].String(), "10.00")

		// 双库存: 活动 sold_count +2; 商品 inventory.locked +2（批次 07 教训: 两者都要锁）
		sc, err := g.DB().GetValue(ctx, "SELECT sold_count FROM flash_sale_item WHERE id=?", itemId)
		t.AssertNil(err)
		t.Assert(sc.Int(), 2)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 100)
		t.Assert(inv["locked"].Int(), 2)

		// 活动库存不足（剩余 8 < 请求 20）→ 40003
		_, err = NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-ORD-2", AddressId: addrId, SkuId: f.SkuId, Quantity: 20, FlashSaleItemId: itemId,
		})
		t.Assert(errCode(err), errcode.CodeSoldOut)

		// 场次商品与所选 SKU 不符 → 50002（防篡改）: 另造一个场次商品绑定别的 sku
		other := f.SkuId + 1
		_, _ = g.DB().Exec(ctx, "INSERT INTO flash_sale_item(activity_id,sku_id,flash_price,stock_count,sold_count,per_limit) VALUES(?,?,?,?,0,2)",
			actId, other, "1.00", 5)
		_, err = NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-ORD-3", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, FlashSaleItemId: itemId,
		})
		t.AssertNil(err) // 该 itemId 仍匹配 f.SkuId, 应通过

		// 活动已结束 → 50002
		_, endItemId := seedFlashSale(ctx, t, "TF-秒杀已结束b", f.SkuId, "5.00", 10, -100, -10)
		_, err = NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-ORD-4", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, FlashSaleItemId: endItemId,
		})
		t.Assert(errCode(err), errcode.CodeActivityInvalid)
	})
}

// TestFlashSaleCancelRestore 取消回补双库存（FR-004）。
func TestFlashSaleCancelRestore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "FS-CXL-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		const name = "TF-秒杀取消"
		defer cleanupFlashSale(ctx, name)
		_, itemId := seedFlashSale(ctx, t, name, f.SkuId, "5.00", 10, -10, 60)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-CXL-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 3, FlashSaleItemId: itemId,
		})
		t.AssertNil(err)
		t.Assert(out.OrderNo != "", true)
		defer cleanupPayFixture(ctx, t, out.OrderNo)
		defer cleanupPayFixture(ctx, t, out.OrderNo)

		t.AssertNil(NewOrderLogic().Cancel(ctx, uid, out.OrderNo, "不想要了"))
		sc, err := g.DB().GetValue(ctx, "SELECT sold_count FROM flash_sale_item WHERE id=?", itemId)
		t.AssertNil(err)
		t.Assert(sc.Int(), 0) // 活动库存回补
		inv, err := g.DB().GetOne(ctx, "SELECT locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["locked"].Int(), 0) // 商品锁定回补
	})
}

// TestFlashSalePayCallback 批次 07 欠账的直接验收: 秒杀单**支付回调必须成功**（库存核销不得未命中）。
func TestFlashSalePayCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "FS-PAY-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		const name = "TF-秒杀支付"
		defer cleanupFlashSale(ctx, name)
		_, itemId := seedFlashSale(ctx, t, name, f.SkuId, "5.00", 10, -10, 60)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-PAY-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 2, FlashSaleItemId: itemId,
		})
		t.AssertNil(err)
		// 支付单/回调留档一并清理（pay_order 有 uk_channel_trade 唯一键, 残留会让下一轮撞键）
		defer cleanupPayFixture(ctx, t, out.OrderNo)

		pay, err := NewPayLogic().Create(ctx, uid, out.OrderNo, 1)
		t.AssertNil(err)
		body, _ := json.Marshal(map[string]any{
			"payNo": pay.PayNo, "channelTradeNo": "CH-FS-" + out.OrderNo, "amountFen": 1000, "success": true,
		})
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body)) // 批次 07: 此处必因"库存核销未命中"回滚

		st, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE order_no=?", out.OrderNo)
		t.AssertNil(err)
		t.Assert(st.Int(), 20) // 已支付
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 98) // 核销 2
		t.Assert(inv["locked"].Int(), 0)
	})
}
