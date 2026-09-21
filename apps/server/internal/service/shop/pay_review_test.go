// pay_review_test.go 支付/退款回调的评审修复回归测试（012 修复轮）。
// 覆盖: C1 退款回调匹配键与状态机 / C2 渠道失败标志 / C3a 单待支付单不变式 /
// C3b 脏态资金标记 / I1 库存核销判行数 / I2 留档出事务与 JSON 列兜底 / I3 expire_time。
// 这些修复此前全部无测试（评审项以"代码注释"形式存在, 无任何断言兜底）。
package shop

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// TestPayCreateClosesPendingPayOrder 单待支付单不变式（C3a）+ 支付截止时间（I3）。
// 原实现每次 Create 都插一张新单, 同订单可并存多张待支付单（重复扣款面）。
func TestPayCreateClosesPendingPayOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-C3-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		_, _ = seedPayFixture(ctx, t, uid, "T-PAY-C3", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-C3")

		first, err := NewPayLogic().Create(ctx, uid, "T-PAY-C3", 1)
		t.AssertNil(err)
		second, err := NewPayLogic().Create(ctx, uid, "T-PAY-C3", 1)
		t.AssertNil(err)
		t.Assert(second.PayNo != first.PayNo, true)

		// 旧单关闭: 90 已关闭（不是 30 支付失败）+ 关闭时间落库
		old, err := g.DB().GetOne(ctx, "SELECT status, closed_time FROM pay_order WHERE pay_no=?", first.PayNo)
		t.AssertNil(err)
		t.Assert(old["status"].Int(), 90)
		t.Assert(old["closed_time"].GTime() != nil, true)

		// 待支付单恰好一张
		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM pay_order WHERE order_no='T-PAY-C3' AND status=10")
		t.AssertNil(err)
		t.Assert(n.Int(), 1)

		// 新单带 30 分钟截止时间（原实现 expire_time 从不写入 → 关单扫描无据可依）
		nw, err := g.DB().GetOne(ctx, "SELECT status, expire_time FROM pay_order WHERE pay_no=?", second.PayNo)
		t.AssertNil(err)
		t.Assert(nw["status"].Int(), 10)
		t.Assert(nw["expire_time"].GTime() != nil, true)
	})
}

// TestPayNotifyRejectsFailFlag 渠道标记支付失败（C2）: 必须拒绝且不进成功路径。
// 原实现完全不看 success 标志 → 一条 success:false 的通知即可把订单刷成已支付（资损面）。
func TestPayNotifyRejectsFailFlag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-C2-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		orderId, skuId := seedPayFixture(ctx, t, uid, "T-PAY-C2", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-C2")
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=2 WHERE sku_id=?", skuId)

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-C2", 1)
		t.AssertNil(err)
		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-C2", "amountFen": 10000, "success": false,
		})
		err = NewPayLogic().HandlePayNotify(ctx, "mock", body)
		t.AssertNE(err, nil)

		pst, err := g.DB().GetValue(ctx, "SELECT status FROM pay_order WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.Assert(pst.Int(), 10)
		ost, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(ost.Int(), 10)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 10)
		t.Assert(inv["locked"].Int(), 2)
		// 留档（未处理）
		ps, err := g.DB().GetValue(ctx,
			"SELECT process_status FROM pay_callback_log WHERE pay_no=? AND notify_type=1", out.PayNo)
		t.AssertNil(err)
		t.Assert(ps.Int(), 0)
	})
}

// TestPayNotifyInventoryMissRollsBack 库存核销未命中（I1）: 必须整体回滚。
// 锁定不足或库存行缺失时原实现静默跳过 → 订单已推进而库存未扣（账实不符）。
func TestPayNotifyInventoryMissRollsBack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-I1-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		orderId, skuId := seedPayFixture(ctx, t, uid, "T-PAY-I1", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-I1")

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-I1", 1)
		t.AssertNil(err)
		// 核销前把锁定清零（模拟锁定不足/库存行异常）
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=0 WHERE sku_id=?", skuId)

		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-I1", "amountFen": 10000, "success": true,
		})
		err = NewPayLogic().HandlePayNotify(ctx, "mock", body)
		t.AssertNE(err, nil)

		// 事务回滚: 支付单与订单都停在原态
		pst, err := g.DB().GetValue(ctx, "SELECT status FROM pay_order WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.Assert(pst.Int(), 10)
		ost, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(ost.Int(), 10)
		// 留档仍要落（成败均留, 事务外写入）
		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM pay_callback_log WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.Assert(n.Int(), 1)
	})
}

// TestPayNotifyDirtyOrderMarker 脏态资金标记（C3b）: 钱已入账而订单不可推进（已取消）。
// 原实现只留一条写在事务内、错误被吞的留档 → 无人知道这笔钱要退。
func TestPayNotifyDirtyOrderMarker(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-C3B-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		orderId, skuId := seedPayFixture(ctx, t, uid, "T-PAY-C3B", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-C3B")
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=2 WHERE sku_id=?", skuId)

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-C3B", 1)
		t.AssertNil(err)
		// 支付发起后订单被取消（已取消订单不可推进）
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET status=90 WHERE id=?", orderId)

		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-C3B", "amountFen": 10000, "success": true,
		})
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body))

		// 支付单保持 20（钱确实到了, 不得篡改事实）; 订单停在 90
		pst, err := g.DB().GetValue(ctx, "SELECT status FROM pay_order WHERE pay_no=?", out.PayNo)
		t.AssertNil(err)
		t.Assert(pst.Int(), 20)
		ost, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(ost.Int(), 90)
		// 对账标记: 留档 process_status=0（未处理）——对账按"资金已收且订单已取消"捞这行
		ps, err := g.DB().GetValue(ctx,
			"SELECT process_status FROM pay_callback_log WHERE pay_no=? AND notify_type=1", out.PayNo)
		t.AssertNil(err)
		t.Assert(ps.Int(), 0)
		// 库存未核销（未推进即不扣库存）
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 10)
		t.Assert(inv["locked"].Int(), 2)
	})
}

// TestPayNotifyDuplicateStillLogged 重复回调也必须留档（I2）。
// 表注释明写"只追加, 重复通知也留档"——重复通知是对账识别重复扣款/刷单的唯一线索。
// 记账更正: 该断言在**批次原始实现**下也通过（当时幂等分支确实写了一条留档, 只是写在事务内）。
// 它真正守住的是修复轮的中间态——半途改动删掉了那行调用并留言"留档由事务外统一落",
// 而当时的代码里并不存在"事务外统一落"（本函数新增的才是）, 于是重复回调一度零留档。
func TestPayNotifyDuplicateStillLogged(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-I2-2"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		_, skuId := seedPayFixture(ctx, t, uid, "T-PAY-I2D", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-I2D")
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=2 WHERE sku_id=?", skuId)

		out, err := NewPayLogic().Create(ctx, uid, "T-PAY-I2D", 1)
		t.AssertNil(err)
		body, _ := json.Marshal(map[string]any{
			"payNo": out.PayNo, "channelTradeNo": "CH-I2D", "amountFen": 10000, "success": true,
		})
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body))
		// 重复回调: 幂等应答, 但必须留档
		t.AssertNil(NewPayLogic().HandlePayNotify(ctx, "mock", body))

		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM pay_callback_log WHERE pay_no=? AND notify_type=1", out.PayNo)
		t.AssertNil(err)
		t.Assert(n.Int(), 2)
		// 第一次已处理(1) + 第二次未处理(0)
		ps, err := g.DB().GetValue(ctx,
			"SELECT SUM(process_status) FROM pay_callback_log WHERE pay_no=? AND notify_type=1", out.PayNo)
		t.AssertNil(err)
		t.Assert(ps.Int(), 1)
	})
}

// TestPayNotifyClosedPayOrder 已关闭支付单收到成功回调（评审 Critical）:
// 必须**不得**当作"幂等已处理"应答 SUCCESS —— 渠道是就旧单扣的款, 钱真的到了。
// 场景（C3a 主动制造）: Create 两次 → 第一张被置 90([已关闭]), 用户在浏览器/微信里扫的仍是第一张码。
// 修复前行为: affected=0 → 直接 return nil → 渠道收到 SUCCESS、支付单停在 90、订单停在 10、
// 库存未核销、**无任何告警**（且 C3b 的对账口径"pay=20 且 order=90"也捞不到它）→ 钱付了、单没了、无人知道。
func TestPayNotifyClosedPayOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-CLOSE-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		orderId, skuId := seedPayFixture(ctx, t, uid, "T-PAY-CLOSED", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-CLOSED")
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=10, locked=2 WHERE sku_id=?", skuId)

		first, err := NewPayLogic().Create(ctx, uid, "T-PAY-CLOSED", 1)
		t.AssertNil(err)
		if _, err = NewPayLogic().Create(ctx, uid, "T-PAY-CLOSED", 1); err != nil { // 第二张把第一张置 90
			t.Fatal(err)
		}

		body, _ := json.Marshal(map[string]any{
			"payNo": first.PayNo, "channelTradeNo": "CH-CLOSE", "amountFen": 10000, "success": true,
		})
		err = NewPayLogic().HandlePayNotify(ctx, "mock", body)
		t.AssertNE(err, nil) // 必须报错（渠道应答 FAIL）, 不得静默 SUCCESS

		// 支付单停在 90; 订单停在 10; 库存未核销
		paySt, err := g.DB().GetValue(ctx, "SELECT status FROM pay_order WHERE pay_no=?", first.PayNo)
		t.AssertNil(err)
		t.Assert(paySt.Int(), 90)
		orderSt, err := g.DB().GetValue(ctx, "SELECT status FROM trade_order WHERE id=?", orderId)
		t.AssertNil(err)
		t.Assert(orderSt.Int(), 10)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 10)
		t.Assert(inv["locked"].Int(), 2)
		// 对账标记: 留档 process_status=0（未处理）——机器可查
		ps, err := g.DB().GetValue(ctx,
			"SELECT process_status FROM pay_callback_log WHERE pay_no=? AND notify_type=1", first.PayNo)
		t.AssertNil(err)
		t.Assert(ps.Int(), 0)
	})
}

// TestPayNotifyMalformedBodyLogged 非法报文留档（I2）: raw_body 为 JSON NOT NULL,
// 非 JSON 原文直接入库报 3140 且被吞 → 最该留档的场景反而零记录。
func TestPayNotifyMalformedBodyLogged(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PAY-I2-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		_, _ = seedPayFixture(ctx, t, uid, "T-PAY-I2", 10000)
		defer cleanupPayFixture(ctx, t, "T-PAY-I2")

		before := maxCallbackLogId(ctx)
		defer cleanupCallbackLogsAfter(ctx, before)
		const bad = "{not-json-payload"
		err := NewPayLogic().HandlePayNotify(ctx, "mock", []byte(bad))
		t.AssertNE(err, nil)

		raw, err := g.DB().GetValue(ctx,
			"SELECT raw_body FROM pay_callback_log WHERE pay_no='' AND notify_type=1 ORDER BY id DESC LIMIT 1")
		t.AssertNil(err)
		t.Assert(strings.Contains(raw.String(), "not-json-payload"), true)
	})
}

// TestRefundNotifyStateMachine 退款回调（C1）: out_refund_no = after_sale_no; 40退款中 → 50已完成;
// affected=0（未匹配/状态不符）必须报错而非记"已处理"。
func TestRefundNotifyStateMachine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			asOK   = "T-AS-OK"
			asWait = "T-AS-WAIT"
			asNone = "T-AS-NONE"
		)
		logStart := maxCallbackLogId(ctx)
		defer cleanupCallbackLogsAfter(ctx, logStart) // 匿名留档（渠道标记失败那条）按 id 窗口清
		cleanup := func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no IN (?,?,?)", asOK, asWait, asNone)
			_, _ = g.DB().Exec(ctx, "DELETE FROM pay_callback_log WHERE pay_no LIKE 'PAY-T-AS%'")
		}
		cleanup()
		defer cleanup()

		seedAfterSale := func(no string, status int) {
			_, err := g.DB().Exec(ctx,
				"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,reason,refund_amount) "+
					"VALUES(?,1,'T-ORD-AS',1,1,1,?,'测试',10.00)", no, status)
			t.AssertNil(err)
		}
		notify := func(payNo, outRefundNo string, success bool) error {
			body, _ := json.Marshal(map[string]any{
				"payNo": payNo, "outRefundNo": outRefundNo, "success": success,
			})
			return NewPayLogic().HandleRefundNotify(ctx, "mock", body)
		}

		// 退款中（40）→ 已完成（50）+ 退款时间落库（原实现查 refund_no 且 30→40, 永不命中）
		seedAfterSale(asOK, 40)
		t.AssertNil(notify("PAY-T-AS-1", asOK, true))
		rec, err := g.DB().GetOne(ctx,
			"SELECT status, refund_time FROM after_sale_order WHERE after_sale_no=?", asOK)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 50)
		t.Assert(rec["refund_time"].GTime() != nil, true)
		ps, err := g.DB().GetValue(ctx,
			"SELECT process_status FROM pay_callback_log WHERE pay_no='PAY-T-AS-1' AND notify_type=2")
		t.AssertNil(err)
		t.Assert(ps.Int(), 1)

		// 状态不符（30待退款, 非退款中）→ 错误 + 留档未处理（原实现 affected=0 也记"已处理"）
		seedAfterSale(asWait, 30)
		err = notify("PAY-T-AS-2", asWait, true)
		t.Assert(errCode(err), errcode.CodeNotFound)
		ps, err = g.DB().GetValue(ctx,
			"SELECT process_status FROM pay_callback_log WHERE pay_no='PAY-T-AS-2' AND notify_type=2")
		t.AssertNil(err)
		t.Assert(ps.Int(), 0)

		// 渠道标记退款失败 → 拒绝
		err = notify("PAY-T-AS-3", asOK, false)
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}
