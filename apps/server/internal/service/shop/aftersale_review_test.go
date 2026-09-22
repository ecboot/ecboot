// aftersale_review_test.go 批次 07 独立评审的回归守卫（每条对应一个已修的 Critical/Important）。
// 评审指出：原 21 个测试全是顺序路径，对"条件更新不重入"这条铁律**零覆盖**——变异实验（去掉
// tryRefund/Cancel 的条件 WHERE、去掉 refund_status 只统计 50 的过滤）后 21 个测试仍全绿。
// 本文件用并发与"在途窗口"把那些不变量变成可断言的。
package shop

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/library/paychannel"
	"ecboot/internal/model"
)

// TestAfterSaleApplyConcurrent 并发申请不得超额（评审 C1）: 行锁事务把"读额度→插入"串行化。
// 修复前无锁无事务: 并发各读一次"额度未占满" → 落多张超额单（评审实证同一行可落 3 张全额单 = 3 倍出款）。
func TestAfterSaleApplyConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-CC-1", "T-AS-CC1", 980000061, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		// 多轮: 单轮的并发窗口未必命中（评审用 3 并发 × 40 轮才 40/40 复现），故反复施加压力
		const rounds, n = 12, 4
		for round := 0; round < rounds; round++ {
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE order_item_id=?", f.ItemId)

			var wg sync.WaitGroup
			errs := make([]error, n)
			start := make(chan struct{})
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					<-start
					_, errs[idx] = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
						OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "并发申请",
					})
				}(i)
			}
			close(start)
			wg.Wait()

			ok := 0
			for _, e := range errs {
				if e == nil {
					ok++
					continue
				}
				// 评审 M-5: 失败必须是**业务码 40008**（配额拒绝）——若出现 10002（如 1205 锁等待超时）
				// 说明锁范围/死锁出了问题, 不能与"正常拒绝"混为一谈（I-2 就是被这样掩盖的）。
				t.Assert(errCode(e), errcode.CodeAfterSaleDenied)
			}
			// 行数量 1 → 每轮恰好 1 张（不让并发穿透"读额度→插入"）
			t.Assert(ok, 1)
			cnt, err := g.DB().GetValue(ctx,
				"SELECT COUNT(*) FROM after_sale_order WHERE order_item_id=?", f.ItemId)
			t.AssertNil(err)
			t.Assert(cnt.Int(), 1)
			sum, err := g.DB().GetValue(ctx,
				"SELECT COALESCE(SUM(refund_amount),0) FROM after_sale_order WHERE order_item_id=?", f.ItemId)
			t.AssertNil(err)
			t.Assert(sum.Float64() <= 10.0, true) // 累计退款额不超行实付
		}
	})
}

// blockingChannel 渠道桩: Refund 阻塞到放行, 用于把"渠道调用中"的窗口确定性地拉长。
type blockingChannel struct {
	entered chan struct{}
	release chan struct{}
}

func (b *blockingChannel) Name() string { return "blocking" }
func (b *blockingChannel) CreateParams(string, int64, string) (map[string]any, error) {
	return nil, nil
}
func (b *blockingChannel) ParseNotify([]byte) (*paychannel.NotifyPayload, error) { return nil, nil }
func (b *blockingChannel) Refund(string, int64) error {
	b.entered <- struct{}{}
	<-b.release
	return nil
}

// TestAfterSaleRefundInFlightNotCancelable 在途退款期间不可撤销（评审 C2）:
// 顺序必须是"先条件推进 30→40, 再调渠道"——渠道未返回时状态已是 40, 买家撤销必须失败。
// 修复前先调渠道再改状态: 该窗口内状态仍是 30（可撤）→ 撤销成功、额度释放、可再退一次（钱已出、单已撤销）。
func TestAfterSaleRefundInFlightNotCancelable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-RACE-1", "T-AS-RACE1", 980000062, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		defer cleanupAfterSaleCallbackLogs(ctx)
		seedAfterSaleRow2(ctx, t, f, "AS-RACE-1", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")

		bc := &blockingChannel{entered: make(chan struct{}, 1), release: make(chan struct{})}
		old := afterSaleRefundChannel
		afterSaleRefundChannel = bc
		defer func() { afterSaleRefundChannel = old }()

		done := make(chan error, 1)
		go func() { done <- NewAfterSaleLogic().RetryRefund(ctx, "AS-RACE-1", "admin:31") }()
		<-bc.entered // 渠道调用中

		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding) // 关键: 状态已推进（先推进再调用）
		t.Assert(errCode(NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-RACE-1")), errcode.CodeStatusNotAllowed)

		bc.release <- struct{}{}
		t.AssertNil(<-done)
	})
}

// TestAfterSaleRefundFailKeepsRefunding 渠道失败后仍不可撤（复审 C2′）:
// 一旦调用过渠道, 失败也必须留在"退款中"（可撤销集 {10,20,30} 之外）。若回退到 30（可撤）:
// 买家撤销 → 可退数量释放 → 以**新幂等键**重新申请并再出款一份（"钱可能已出"变成可撤状态 = 新超额退款面,
// 复审已单线程复现）。本用例钉住"渠道失败 ⇒ 状态仍 40 ⇒ 不可撤 ⇒ 额度不释放"。
func TestAfterSaleRefundFailKeepsRefunding(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-C2B-1", "T-AS-C2B1", 980000071, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-C2B-1", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")

		spy := &refundSpy{failWith: errors.New("渠道超时")}
		defer useRefundSpy(spy)()

		// 30 → 发起（失败: 模糊失败）→ 必须留在 40
		t.Assert(errCode(NewAfterSaleLogic().RetryRefund(ctx, "AS-C2B-1", "admin:41")), errcode.CodeRefundFailed)
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding)
		t.Assert(rec["fail_reason"].String() != "", true)

		// 关键: 此时撤销必须失败（否则额度释放 → 新键再退一份）
		t.Assert(errCode(NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-C2B-1")), errcode.CodeStatusNotAllowed)

		// 额度未被释放: 再申请同数量仍被拒（累计口径）
		_, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "再申请",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)
	})
}

// TestAfterSaleRetryRefundConcurrent 并发重试只推进一次（评审 I5 变异①的守卫）。
func TestAfterSaleRetryRefundConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-RTY-C", "T-AS-RTYC", 980000063, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-RTY-C", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		const n = 6
		var wg sync.WaitGroup
		oks := make([]bool, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				oks[idx] = NewAfterSaleLogic().RetryRefund(ctx, "AS-RTY-C", "admin:32") == nil
			}(i)
		}
		close(start)
		wg.Wait()

		success := 0
		for _, b := range oks {
			if b {
				success++
			}
		}
		t.Assert(success, 1) // 只有一次推进成功（其余 affected=0 → 40006）
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding)
	})
}

// TestAfterSaleCancelConcurrent 并发撤销只生效一次（评审 I5 变异②的守卫: 条件更新 WHERE 的语义）。
func TestAfterSaleCancelConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-CAN-C", "T-AS-CANC", 980000064, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		// 多轮 × 多并发: 条件更新的 WHERE 缺失时，两个都通过前置校验的调用会"双成功"
		const rounds, n = 12, 8
		for round := 0; round < rounds; round++ {
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no='AS-CAN-C'")
			seedAfterSaleRow2(ctx, t, f, "AS-CAN-C", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "10.00")

			var wg sync.WaitGroup
			oks := make([]bool, n)
			start := make(chan struct{})
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					<-start
					oks[idx] = NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-CAN-C") == nil
				}(i)
			}
			close(start)
			wg.Wait()

			success := 0
			for _, b := range oks {
				if b {
					success++
				}
			}
			t.Assert(success, 1) // 每轮恰好一次生效
		}
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleCanceled)
	})
}

// TestAfterSaleRefundStatusOnlyFinished 订单退款状态只统计"已完成"的申请（评审 I5 变异③的守卫）:
// 在途(40)/已撤销(91)的申请不得让订单提前变成"部分/全额退款"。
func TestAfterSaleRefundStatusOnlyFinished(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-RS-1", "T-AS-RS1", 980000065, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		// 一行数量 2: 已完成 1 件 + 在途退款中 1 件 + 已撤销 1 件
		seedAfterSaleRow2(ctx, t, f, "AS-RS-DONE", afterSaleFinished, afterSaleTypeReturnGoods, 1, "10.00")
		seedAfterSaleRow2(ctx, t, f, "AS-RS-FLIGHT", afterSaleRefunding, afterSaleTypeRefundOnly, 1, "10.00")
		seedAfterSaleRow2(ctx, t, f, "AS-RS-CANCEL", afterSaleCanceled, afterSaleTypeRefundOnly, 1, "10.00")

		t.AssertNil(recalcOrderRefundStatus(ctx, f.OrderId))
		st, err := g.DB().GetValue(ctx, "SELECT refund_status FROM trade_order WHERE id=?", f.OrderId)
		t.AssertNil(err)
		// 只看已完成: 1/2 → 部分退款（1）; 若把在途也算进去会变成 2（全额）
		t.Assert(st.Int(), 1)
	})
}

// TestAfterSaleCancelLostRace 撤销**输掉**状态迁移时不得覆盖状态（复审给出的零测试缝黑盒构造）。
// 构造: 外部事务先在**未提交**状态下把行改成终态（持该行 X 锁）→ 目标方法的前置**普通读**仍见旧状态
// （InnoDB MVCC 看不到未提交改动）→ 其条件 UPDATE 阻塞在锁上 → 外部事务提交后条件重新求值 → 不命中
// → affected=0 → 必须返 40006 且**不得覆盖状态**。去掉条件 WHERE 的变异会让本用例红（复审已实证）。
func TestAfterSaleCancelLostRace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-LOST-C", "T-AS-LOSTC", 980000081, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-LOST-C", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "10.00")

		tx, err := g.DB().Begin(ctx)
		t.AssertNil(err)
		_, err = tx.Exec("UPDATE after_sale_order SET status=? WHERE after_sale_no=?", afterSaleCanceled, "AS-LOST-C")
		t.AssertNil(err)

		// 池预热: 确保障碍物之外的连接已就绪, 否则 goroutine 可能因取不到连接而"迟到",
		// 其前置读会落在外部提交之后（实测会导致本构造无法区分条件 WHERE 的变异）。
		var warm sync.WaitGroup
		for i := 0; i < 4; i++ {
			warm.Add(1)
			go func() { defer warm.Done(); _, _ = g.DB().GetValue(ctx, "SELECT 1") }()
		}
		warm.Wait()

		done := make(chan error, 1)
		go func() { done <- NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-LOST-C") }()
		time.Sleep(2000 * time.Millisecond) // 让 Cancel 通过前置校验并阻塞在条件 UPDATE 上
		t.AssertNil(tx.Commit())

		t.Assert(errCode(<-done), errcode.CodeStatusNotAllowed)
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleCanceled) // 终态未被覆盖
	})
}

// TestAfterSaleRetryRefundLostRace 重试退款**未赢得状态迁移就不得出款**（复审给出的黑盒构造）:
// 外部事务未提交地把行改成终态 → RetryRefund 的条件推进阻塞 → 提交后不命中 → 40006 且**渠道调用 0 次**。
func TestAfterSaleRetryRefundLostRace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-LOST-R", "T-AS-LOSTR", 980000082, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-LOST-R", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		tx, err := g.DB().Begin(ctx)
		t.AssertNil(err)
		_, err = tx.Exec("UPDATE after_sale_order SET status=? WHERE after_sale_no=?", afterSaleCanceled, "AS-LOST-R")
		t.AssertNil(err)

		// 池预热（同 CancelLostRace 的理由）: 避免 goroutine 因取连接而迟到,
		// 使其前置读落在外部提交之后, 那样构造就失去了区分力。
		var warm sync.WaitGroup
		for i := 0; i < 4; i++ {
			warm.Add(1)
			go func() { defer warm.Done(); _, _ = g.DB().GetValue(ctx, "SELECT 1") }()
		}
		warm.Wait()

		done := make(chan error, 1)
		go func() { done <- NewAfterSaleLogic().RetryRefund(ctx, "AS-LOST-R", "admin:51") }()
		time.Sleep(2000 * time.Millisecond)
		t.AssertNil(tx.Commit())

		t.Assert(errCode(<-done), errcode.CodeStatusNotAllowed)
		t.Assert(spy.calls, 0) // 关键: 没赢得状态迁移, 就绝不能调用渠道
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleCanceled)
	})
}
