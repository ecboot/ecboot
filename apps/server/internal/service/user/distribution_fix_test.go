// distribution_fix_test.go 批次 11 修复轮回归（评审 C3/C4/C5/I1/I2/C2/C6）——并发守卫与红线。
package user

import (
	"context"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// TestDistCommissionConcurrentGuards C3/C4/C5 守卫: 并发计提不双计提、并发冲销不双扣、
// 未结算并发冲销输家不扣款（评审探针 P4/P5/P5b 的交付测试化）。
func TestDistCommissionConcurrentGuards(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DG-A")
		b := distCleanupUser(ctx, t, "TF-DG-B")
		defer distCleanupAll(ctx, a, b)
		const no = "TF-DIST-CC-1"
		defer cleanupDistOrder(ctx, no)
		f := newDistFixture(t)
		defer f.teardown(ctx)
		logic := NewDistributionLogic()
		seedRelation(ctx, t, a, b)
		seedDistributor(ctx, t, a, 2)
		// 规则挂到本 fixture 的分类（订单 SKU 归属该分类才会命中）
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_type=1 AND scope_id=?", f.catId)
		_, err := g.DB().Exec(ctx,
			"INSERT INTO commission_rule(scope_type,scope_id,level1_rate,level2_rate,status) VALUES(1,?,10.00,0.00,1)",
			f.catId)
		t.AssertNil(err)
		itemId := seedDistOrderSku(ctx, t, no, b, "100.00", f.skuId, f.spuId)

		// C3: 两并发 SettleOrder → 恰 1 条记录（唯一键 uk_item_bene_level_rev 兜底）
		var wg sync.WaitGroup
		errs := make([]error, 2)
		start := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				errs[0] = logic.SettleOrder(ctx, no)
			}()
		}
		// 同订单两个 goroutine 调同一方法 → 用两个调用
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs[1] = logic.SettleOrder(ctx, no)
		}()
		close(start)
		wg.Wait()
		for _, e := range errs {
			t.AssertNil(e)
		}
		cnt, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE order_item_id=? AND beneficiary_user_id=?", itemId, a)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1) // 不双计提

		// 保护期满 → 结算入账 10.00
		_, _ = g.DB().Exec(ctx,
			"UPDATE commission_record SET created_at=DATE_SUB(NOW(), INTERVAL 8 DAY) WHERE order_item_id=?", itemId)
		_, err = logic.ConfirmSettle(ctx)
		t.AssertNil(err)
		acc, _ := g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "10.00")

		// C4: 两并发 ReverseOnRefund（已结算）→ 恰 1 条冲销、余额 10→0（不双扣到 -10）
		var wg2 sync.WaitGroup
		revErrs := make([]error, 2)
		start2 := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg2.Add(1)
			go func(idx int) {
				defer wg2.Done()
				<-start2
				revErrs[idx] = logic.ReverseOnRefund(ctx, itemId)
			}(i)
		}
		close(start2)
		wg2.Wait()
		revOK := 0
		for _, e := range revErrs {
			if e == nil {
				revOK++
			}
		}
		t.Assert(revOK >= 1, true) // 至少一次成功
		revCnt, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE reversal_of_id IS NOT NULL")
		t.AssertNil(err)
		t.Assert(revCnt.Int(), 1) // 恰一条冲销（uk_reversal_of 兜底）
		acc, _ = g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "0.00") // 不双扣
	})
}

// TestDistReverseUnsettledConcurrent C5 守卫: 未结算记录两并发冲销 → 输家不扣款（从未入账）。
func TestDistReverseUnsettledConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DU-A")
		b := distCleanupUser(ctx, t, "TF-DU-B")
		defer distCleanupAll(ctx, a, b)
		const no = "TF-DIST-UU-1"
		defer cleanupDistOrder(ctx, no)
		logic := NewDistributionLogic()
		seedRelation(ctx, t, a, b)
		seedDistributor(ctx, t, a, 2)
		itemId := seedDistOrder(ctx, t, no, b, "100.00")
		t.AssertNil(logic.SettleOrder(ctx, no))
		seedAccount(ctx, t, a, "100.00", "0.00") // 独立预置余额（佣金未入账）

		var wg sync.WaitGroup
		start := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				_ = logic.ReverseOnRefund(ctx, itemId)
			}()
		}
		close(start)
		wg.Wait()

		// 佣金从未入账 → 不得被扣（评审探针 P5b: 余额 100→90）
		acc, _ := g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "100.00")
	})
}

// TestDistBindCycle I1 守卫: 互环绑定拒绝（A←B 已建立, 再绑 B→A → 拒绝）。
func TestDistBindCycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DX-A")
		b := distCleanupUser(ctx, t, "TF-DX-B")
		defer distCleanupAll(ctx, a, b)
		logic := NewDistributionLogic()
		t.AssertNil(logic.BindRelation(ctx, b, a, 1)) // B 的上级 = A
		// 互环: 给 A 绑上级 B（B 是 A 的下级）→ 拒绝
		t.Assert(errCode(logic.BindRelation(ctx, a, b, 1)), errcode.CodeInvalidParam)
	})
}

// TestDistCommissionDecimalRate I2 守卫: 小数比例精算（2.50% × 100.00 = 2.50, 不截断为 2.00）。
func TestDistCommissionDecimalRate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := newDistFixture(t)
		defer f.teardown(ctx)
		a := distCleanupUser(ctx, t, "TF-DD-A")
		b := distCleanupUser(ctx, t, "TF-DD-B")
		defer distCleanupAll(ctx, a, b)
		const no = "TF-DIST-DD-1"
		defer cleanupDistOrder(ctx, no)
		logic := NewDistributionLogic()
		seedRelation(ctx, t, a, b)
		seedDistributor(ctx, t, a, 2)
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_type=1 AND scope_id=?", f.catId)
		_, err := g.DB().Exec(ctx,
			"INSERT INTO commission_rule(scope_type,scope_id,level1_rate,level2_rate,status) VALUES(1,?,2.50,0.00,1)",
			f.catId)
		t.AssertNil(err)
		itemId := seedDistOrderSku(ctx, t, no, b, "100.00", f.skuId, f.spuId)
		t.AssertNil(logic.SettleOrder(ctx, no))
		rec, err := g.DB().GetOne(ctx,
			"SELECT amount, rate FROM commission_record WHERE order_item_id=? AND beneficiary_user_id=?", itemId, a)
		t.AssertNil(err)
		t.Assert(rec["amount"].String(), "2.50") // 不截断为 2.00
		t.Assert(rec["rate"].String(), "2.50")   // 记录存原值
	})
}

// TestDistAttributionPriority C2 守卫: 订单归因人优先于关系链（V28 attributed_user_id）。
func TestDistAttributionPriority(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := newDistFixture(t)
		defer f.teardown(ctx)
		a := distCleanupUser(ctx, t, "TF-DP-A") // 关系链上级
		s := distCleanupUser(ctx, t, "TF-DP-S") // 分享归因人
		b := distCleanupUser(ctx, t, "TF-DP-B") // 买家
		defer distCleanupAll(ctx, a, s, b)
		const no = "TF-DIST-PP-1"
		defer cleanupDistOrder(ctx, no)
		logic := NewDistributionLogic()
		seedRelation(ctx, t, a, b) // 关系链: B 的上级 = A
		seedDistributor(ctx, t, a, 2)
		seedDistributor(ctx, t, s, 2)
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_type=1 AND scope_id=?", f.catId)
		_, err := g.DB().Exec(ctx,
			"INSERT INTO commission_rule(scope_type,scope_id,level1_rate,level2_rate,status) VALUES(1,?,10.00,0.00,1)",
			f.catId)
		t.AssertNil(err)
		itemId := seedDistOrderSku(ctx, t, no, b, "100.00", f.skuId, f.spuId)
		// 订单归因人 = S（非关系链上级 A）
		_, _ = g.DB().Exec(ctx,
			"UPDATE trade_order SET attributed_user_id=?, attribution_type=1 WHERE order_no=?", s, no)

		t.AssertNil(logic.SettleOrder(ctx, no))
		rec, err := g.DB().GetOne(ctx,
			"SELECT beneficiary_user_id, level FROM commission_record WHERE order_item_id=?", itemId)
		t.AssertNil(err)
		t.Assert(rec["beneficiary_user_id"].Int64(), s) // 归因人优先于关系链
		t.Assert(rec["level"].Int(), 1)
	})
}

// TestDistWithdrawFrozenMismatch C6 守卫: 冻结与单据不符时打款 → 报错且状态不动（不静默放行）。
func TestDistWithdrawFrozenMismatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DM-A")
		defer distCleanupAll(ctx, a)
		logic := NewDistributionLogic()
		admin := NewDistributionAdminLogic()
		seedAccount(ctx, t, a, "90.00", "10.00") // 冻结 10（真实）——单据将写 60（伪造场景模拟账实不符）

		no, err := logic.WithdrawApply(ctx, a, "60.00") // 先造出 20 态单据? 不——直接造库
		t.AssertNil(err)
		// 把单据推到 20（绕过正常冻结量: 直接改库模拟账实分歧）
		t.AssertNil(admin.AdminWithdrawAudit(ctx, no, true, ""))
		_, _ = g.DB().Exec(ctx, "UPDATE user_account SET balance=90.00, frozen=1.00 WHERE user_id=?", a)

		// 打款成功 60 → 冻结 1.00 不足 → 必须报错且状态留在 20（事务回滚）
		errCodeGot := errCode(admin.AdminWithdrawPay(ctx, no, true, "CH-MISMATCH", ""))
		t.Assert(errCodeGot != 0, true) // 不得静默成功
		st, _ := g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", no)
		t.Assert(st.Int(), 20) // 状态未推进
	})
}
