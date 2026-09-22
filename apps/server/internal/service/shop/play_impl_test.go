// play_impl_test.go 砍价/助力玩法动作（015-marketing-c 批次 09）。
// 重点: 一人一刀/一人一助力（唯一键兜底）、不砍穿（末刀补齐到恰好底价）、并发只生效一次、
// 助力达标**只投递一次**发奖意图、风控端口拦截。
package shop

import (
	"context"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 风控桩（IRiskHit）----

type riskSpy struct {
	calls   int
	blocked bool
}

func (r *riskSpy) Hit(_ context.Context, _ int64, _ int, _ string) (bool, error) {
	r.calls++
	return r.blocked, nil
}

func useRiskSpy(spy *riskSpy) func() {
	old := RiskHit
	RiskHit = spy
	return func() { RiskHit = old }
}

// ---- 发奖意图桩（IAssistReward）----

type rewardSpy struct {
	calls int
}

func (r *rewardSpy) GrantForAssist(_ context.Context, _, _, _ int64, _ int, _ int64) { r.calls++ }

func useRewardSpy(spy *rewardSpy) func() {
	old := AssistReward
	AssistReward = spy
	return func() { AssistReward = old }
}

// cutErr/helpErr 只取错误（Cut/Help 返回 (结果, 错误) 二元组, 便于断言错误码）。
func cutErr(ctx context.Context, uid, recordId int64) error {
	_, e := NewBargainLogic().Cut(ctx, uid, recordId)
	return e
}

func helpErr(ctx context.Context, uid, recordId int64) error {
	_, e := NewAssistLogic().Help(ctx, uid, recordId)
	return e
}

// ---- fixture ----

// seedBargain 造砍价活动 + 场次商品（original/floor/maxCuts 给定; 时间为进行中）。
func seedBargain(ctx context.Context, t *gtest.T, name string, spuId, skuId int64, original, floor string, maxCuts int) (int64, int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_item WHERE activity_id IN (SELECT id FROM bargain_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_activity WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO bargain_activity(name,spu_id,start_time,end_time,status,deleted) "+
			"VALUES(?,?,DATE_SUB(NOW(), INTERVAL 1 HOUR),DATE_ADD(NOW(), INTERVAL 1 HOUR),1,0)", name, spuId)
	t.AssertNil(err)
	actId, _ := res.LastInsertId()
	ires, err := g.DB().Exec(ctx,
		"INSERT INTO bargain_item(activity_id,sku_id,original_price,floor_price,max_cut_count,config) VALUES(?,?,?,?,?,'{}')",
		actId, skuId, original, floor, maxCuts)
	t.AssertNil(err)
	itemId, _ := ires.LastInsertId()
	return actId, itemId
}

func cleanupBargain(ctx context.Context, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_helper WHERE record_id IN (SELECT id FROM bargain_record WHERE item_id IN (SELECT id FROM bargain_item WHERE activity_id IN (SELECT id FROM bargain_activity WHERE name=?)))", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_record WHERE item_id IN (SELECT id FROM bargain_item WHERE activity_id IN (SELECT id FROM bargain_activity WHERE name=?))", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_item WHERE activity_id IN (SELECT id FROM bargain_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_activity WHERE name=?", name)
}

// seedAssist 造助力活动（需 required 人 / 每人限 perLimit 次）。
func seedAssist(ctx context.Context, t *gtest.T, name string, required, perLimit, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_helper WHERE record_id IN (SELECT id FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?))", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_activity WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO assist_activity(name,reward_type,reward_ref,required_count,per_limit,start_time,end_time,status,deleted) "+
			"VALUES(?,1,1,?,?,DATE_SUB(NOW(), INTERVAL 1 HOUR),DATE_ADD(NOW(), INTERVAL 1 HOUR),?,0)",
		name, required, perLimit, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupAssist(ctx context.Context, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_helper WHERE record_id IN (SELECT id FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?))", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_activity WHERE name=?", name)
}

// ---- US3 砍价 ----

// TestBargainLaunchCut 砍价闭环（FR-007~010）: 发起首刀 → 他人帮砍 → 到底价 → 状态可下单; 三种拒绝。
func TestBargainLaunchCut(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2 = "BGN-1", "BGN-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-砍价闭环"
		defer cleanupBargain(ctx, name)
		// 原价 100.00 / 底价 60.00 / 5 刀 → 每刀 8.00
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 5)

		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		t.Assert(out.CurrentPrice, "92.00") // 发起即首刀

		// 进度可见
		p, err := NewBargainLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.CurrentPrice, "92.00")
		t.Assert(p.OriginalPrice, "100.00")
		t.Assert(p.FloorPrice, "60.00")
		t.Assert(p.CutCount, 1)
		t.Assert(p.Status, 1)
		t.Assert(len(p.Helpers), 1)
		t.Assert(p.Helpers[0].Nickname != "", true) // 脱敏昵称非空

		// 他人帮砍一刀
		cut, err := NewBargainLogic().Cut(ctx, u2, out.RecordId)
		t.AssertNil(err)
		t.Assert(cut.CutAmount, "8.00")
		t.Assert(cut.CurrentPrice, "84.00")
		t.Assert(cut.FloorReached, false)

		// 同一人再砍 → 50004; 发起者本人砍 → 50002
		t.Assert(errCode(cutErr(ctx, u2, out.RecordId)), errcode.CodeAlreadyCut)
		t.Assert(errCode(cutErr(ctx, u1, out.RecordId)), errcode.CodeActivityInvalid)
	})
}

// TestBargainFloorLastCut 末刀补齐到恰好底价（FR-009: 不砍穿、多刀之和 = 起始价−底价）。
func TestBargainFloorLastCut(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2, h3 = "BGF-1", "BGF-2", "BGF-3"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		defer cleanupPointUser(ctx, t, h3)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)
		u3 := seedPointUser(ctx, t, h3)

		const name = "TF-砍价尾差"
		defer cleanupBargain(ctx, name)
		// 原价 100.00 / 底价 60.00 / 3 刀 → 每刀 13.33（4000/3 向下取整）; 末刀补齐 13.34
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 3)

		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		t.Assert(out.CurrentPrice, "86.67") // 100 - 13.33
		c2, err := NewBargainLogic().Cut(ctx, u2, out.RecordId)
		t.AssertNil(err)
		t.Assert(c2.CurrentPrice, "73.34")
		c3, err := NewBargainLogic().Cut(ctx, u3, out.RecordId)
		t.AssertNil(err)
		t.Assert(c3.CurrentPrice, "60.00") // 末刀补齐到恰好底价（不砍穿）
		t.Assert(c3.FloorReached, true)

		p, err := NewBargainLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.Status, 2) // 到底价
		t.Assert(p.CutCount, 3)
	})
}

// TestBargainExpired 超时单不可砍（FR-010）。
func TestBargainExpired(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2 = "BGE-1", "BGE-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-砍价超时"
		defer cleanupBargain(ctx, name)
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 5)
		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		_, _ = g.DB().Exec(ctx, "UPDATE bargain_record SET expire_time=DATE_SUB(NOW(), INTERVAL 1 HOUR) WHERE id=?", out.RecordId)

		t.Assert(errCode(cutErr(ctx, u2, out.RecordId)), errcode.CodeActivityInvalid)
		p, err := NewBargainLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.Status, 4) // 惰性判定为超时
	})
}

// TestBargainCutConcurrent 并发帮砍: 同一人只成功一次（唯一键兜底 + 1062 转业务码）。
func TestBargainCutConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2 = "BGC-1", "BGC-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-砍价并发"
		defer cleanupBargain(ctx, name)
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 10)
		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)

		const n = 6
		var wg sync.WaitGroup
		oks := make([]bool, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				_, e := NewBargainLogic().Cut(ctx, u2, out.RecordId)
				oks[idx] = e == nil
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
		t.Assert(success, 1) // 一人一刀
		p, err := NewBargainLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.CutCount, 2) // 首刀 + 一刀
	})
}

// TestBargainRiskBlocked 帮砍入口风控拦截（FR-015）: 端口返回拦截 → 10005。
func TestBargainRiskBlocked(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2 = "BGR-1", "BGR-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-砍价风控"
		defer cleanupBargain(ctx, name)
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 5)
		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)

		spy := &riskSpy{blocked: true}
		defer useRiskSpy(spy)()
		t.Assert(errCode(cutErr(ctx, u2, out.RecordId)), errcode.CodeForbidden)
		t.Assert(spy.calls, 1)
	})
}

// ---- US4 助力 ----

// TestAssistLaunchHelp 助力闭环（FR-012~014）: 发起 → 三人助力 → 达标发奖意图一次。
func TestAssistLaunchHelp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1, h2, h3, h4 = "AST-1", "AST-2", "AST-3", "AST-4"
		for _, h := range []string{h1, h2, h3, h4} {
			defer cleanupPointUser(ctx, t, h)
		}
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)
		u3 := seedPointUser(ctx, t, h3)
		u4 := seedPointUser(ctx, t, h4)

		const name = "TF-助力闭环"
		defer cleanupAssist(ctx, name)
		actId := seedAssist(ctx, t, name, 3, 1, 1)

		spy := &rewardSpy{}
		defer useRewardSpy(spy)()

		out, err := NewAssistLogic().Launch(ctx, u1, actId)
		t.AssertNil(err)
		p, err := NewAssistLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.HelperCount, 0)
		t.Assert(p.RequiredCount, 3)
		t.Assert(p.Status, 1)

		r2, err := NewAssistLogic().Help(ctx, u2, out.RecordId)
		t.AssertNil(err)
		t.Assert(r2.Done, false)
		r3, err := NewAssistLogic().Help(ctx, u3, out.RecordId)
		t.AssertNil(err)
		t.Assert(r3.Done, false)
		r4, err := NewAssistLogic().Help(ctx, u4, out.RecordId)
		t.AssertNil(err)
		t.Assert(r4.Done, true) // 第三人达标

		p, err = NewAssistLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.HelperCount, 3)
		t.Assert(p.Status, 2) // 已完成发奖
		t.Assert(p.FinishTime != "", true)
		t.Assert(len(p.Helpers), 3)
		t.Assert(spy.calls, 1) // 发奖意图恰好一次
	})
}

// TestAssistRules 助力规则: 重复助力 50007 / 本人助力 50002 / 发起超限 50006 / 活动停用 50002。
func TestAssistRules(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1, h2 = "ASR-1", "ASR-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-助力规则"
		defer cleanupAssist(ctx, name)
		actId := seedAssist(ctx, t, name, 3, 1, 1) // 每人限发起 1 次

		out, err := NewAssistLogic().Launch(ctx, u1, actId)
		t.AssertNil(err)
		// 首次助力成功 → 重复助力 50007
		if _, e := NewAssistLogic().Help(ctx, u2, out.RecordId); e != nil {
			t.Fatal(e)
		}
		_, err = NewAssistLogic().Help(ctx, u2, out.RecordId)
		t.Assert(errCode(err), errcode.CodeAlreadyAssisted)
		// 本人助力自己的记录
		t.Assert(errCode(helpErr(ctx, u1, out.RecordId)), errcode.CodeActivityInvalid)
		// 发起超限（u1 已发起 1 次, 上限 1）
		_, err = NewAssistLogic().Launch(ctx, u1, actId)
		t.Assert(errCode(err), errcode.CodeAssistUsedUp)

		// 停用活动 → 50002
		const name2 = "TF-助力停用"
		defer cleanupAssist(ctx, name2)
		act2 := seedAssist(ctx, t, name2, 2, 1, 0)
		_, err = NewAssistLogic().Launch(ctx, u2, act2)
		t.Assert(errCode(err), errcode.CodeActivityInvalid)
	})
}

// TestAssistHelpConcurrent 并发达标只投递一次发奖意图（FR-014）。
func TestAssistHelpConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		hashes := []string{"ASC-1", "ASC-2", "ASC-3", "ASC-4", "ASC-5", "ASC-6"}
		uids := make([]int64, 0, len(hashes))
		for _, h := range hashes {
			defer cleanupPointUser(ctx, t, h)
			uids = append(uids, seedPointUser(ctx, t, h))
		}
		const name = "TF-助力并发"
		defer cleanupAssist(ctx, name)
		actId := seedAssist(ctx, t, name, 2, 1, 1) // 需 2 人 → 并发助力者超过 2

		spy := &rewardSpy{}
		defer useRewardSpy(spy)()
		out, err := NewAssistLogic().Launch(ctx, uids[0], actId)
		t.AssertNil(err)

		var wg sync.WaitGroup
		start := make(chan struct{})
		for _, uid := range uids[1:] { // 5 人并发助力, 达标只需 2
			wg.Add(1)
			go func(u int64) {
				defer wg.Done()
				<-start
				_, _ = NewAssistLogic().Help(ctx, u, out.RecordId)
			}(uid)
		}
		close(start)
		wg.Wait()

		p, err := NewAssistLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.Status, 2)  // 达标
		t.Assert(spy.calls, 1) // 发奖意图**只一次**
	})
}

// TestAssistRiskBlocked 助力入口风控拦截（FR-015）。
func TestAssistRiskBlocked(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1, h2 = "ASB-1", "ASB-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-助力风控"
		defer cleanupAssist(ctx, name)
		actId := seedAssist(ctx, t, name, 2, 1, 1)
		out, err := NewAssistLogic().Launch(ctx, u1, actId)
		t.AssertNil(err)

		spy := &riskSpy{blocked: true}
		defer useRiskSpy(spy)()
		_, err = NewAssistLogic().Help(ctx, u2, out.RecordId)
		t.Assert(errCode(err), errcode.CodeForbidden)
	})
}

// TestMarketingLists 五类公开列表: 只出进行中 + 首页聚合四块（空块返回空数组）。
func TestMarketingLists(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)

		const bName, aName = "TF-列表砍价", "TF-列表助力"
		defer cleanupBargain(ctx, bName)
		defer cleanupAssist(ctx, aName)
		seedBargain(ctx, t, bName, f.SpuId, f.SkuId, "100.00", "60.00", 5)
		seedAssist(ctx, t, aName, 2, 1, 1)

		bl, err := NewMarketingLogic().PublicBargains(ctx, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		found := false
		for _, it := range bl.List {
			if it.Name == bName {
				found = true
				t.Assert(len(it.Items), 1)
				t.Assert(it.Items[0].FloorPrice, "60.00")
				t.Assert(it.Items[0].MaxCutCount, 5)
			}
		}
		t.Assert(found, true)

		al, err := NewMarketingLogic().PublicAssists(ctx, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		found = false
		for _, it := range al.List {
			if it.Name == aName {
				found = true
				t.Assert(it.RequiredCount, 2)
				t.Assert(it.RewardType, 1)
			}
		}
		t.Assert(found, true)

		// 首页聚合: 四块都返回（券/轮播/楼层可能为空 → 空数组而非 null）
		idx, err := NewMarketingLogic().Index(ctx)
		t.AssertNil(err)
		t.Assert(idx.Banners != nil, true)
		t.Assert(idx.Floors != nil, true)
		t.Assert(idx.Coupons != nil, true)
		t.Assert(idx.Bargains != nil, true)
		hitAssist := false
		for _, it := range idx.Assists {
			if it.Name == aName {
				hitAssist = true
			}
		}
		t.Assert(hitAssist, true)
	})
}
