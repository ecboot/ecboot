// play_impl_test.go 砍价/助力玩法动作（015-marketing-c 批次 09）。
// 重点: 一人一刀/一人一助力（唯一键兜底）、不砍穿（末刀补齐到恰好底价）、并发只生效一次、
// 助力达标**只投递一次**发奖意图、风控端口拦截。
package shop

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/dao"
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

// TestBargainCutConcurrent 并发帮砍（C2 回归, 评审实证场景）: **不同会员**同时帮砍。
// 原实现的守卫只有 status=1（中间刀不改状态→并发全命中）, current_price 被各自陈旧读覆盖
// （评审探针: 实降 8.88 元 vs 流水和 39.96 元）; 同一 uid 并发的场景只走唯一键路径, 测不到它。
// 不变量: Σ(帮砍流水, 含首刀) == 起始价 − 当前价; cut_count 满时必须 status=2 且恰好到底价。
func TestBargainCutConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		// 1 发起人 + 8 名不同帮砍人
		const n = 8
		uids := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			h := fmt.Sprintf("BGC-%d", i)
			defer cleanupPointUser(ctx, t, h)
			uids[i] = seedPointUser(ctx, t, h)
		}

		const name = "TF-砍价并发"
		defer cleanupBargain(ctx, name)
		// 100/60/9 刀 → 每刀 4.44; 发起首刀 4.44 + 8 人帮砍 = 刀数满（9）, 第 8 人补齐恰好到底价
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 9)
		out, err := NewBargainLogic().Launch(ctx, uids[0], itemId)
		t.AssertNil(err)

		var wg sync.WaitGroup
		oks := make([]bool, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				_, e := NewBargainLogic().Cut(ctx, uids[idx+1], out.RecordId)
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
		t.Assert(success, n) // 8 名不同会员**每人都成功一刀**（串行化, 不再丢失更新）

		rec, err := g.DB().GetOne(ctx,
			"SELECT current_price, cut_count, status, success_time FROM bargain_record WHERE id=?", out.RecordId)
		t.AssertNil(err)
		t.Assert(rec["current_price"].String(), "60.00") // 恰好到底价（不砍穿）
		t.Assert(rec["cut_count"].Int(), 9)              // 首刀 + 8
		t.Assert(rec["status"].Int(), 2)                 // 到底价待下单
		t.Assert(rec["success_time"].IsNil(), false)     // success_time 已写（FR-006 校验依据）

		// **守恒不变量**: Σ(帮砍流水, 含首刀) == 起始价 − 底价 == 40.00
		sumV, err := g.DB().GetValue(ctx,
			"SELECT COALESCE(SUM(cut_amount),0) FROM bargain_helper WHERE record_id=?", out.RecordId)
		t.AssertNil(err)
		sum, err := yuanFen(sumV.String())
		t.AssertNil(err)
		t.Assert(sum, 4000)
	})
}

// TestBargainSecondUserLaunch C1 回归: 第二个会员对同一场次商品发起砍价**必须成功**。
// 原缺陷: uk_order_no(order_no) + NOT NULL DEFAULT '' 使全库只能存在一张砍价单。
func TestBargainSecondUserLaunch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2 = "BGS-1", "BGS-2"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)

		const name = "TF-砍价双人发起"
		defer cleanupBargain(ctx, name)
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 5)

		out1, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		out2, err := NewBargainLogic().Launch(ctx, u2, itemId) // 第二人发起: 修复前必 1062
		t.AssertNil(err)
		t.Assert(out1.RecordId != out2.RecordId, true)

		cnt, err := dao.BargainRecord.Ctx(ctx).Where(dao.BargainRecord.Columns().ItemId, itemId).Count()
		t.AssertNil(err)
		t.Assert(cnt, 2) // 两张砍价单共存
	})
}

// TestBargainOrderPlacement I4 回归: 砍价成交价**下单闭环**（spec 006 FR-006）。
// 到底价后按 current_price 成交 → 砍价单置 3 + 回填 order_no; 重复下单/数量不符/SKU 不符 → 40004。
func TestBargainOrderPlacement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2, h3 = "BGO-1", "BGO-2", "BGO-3"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		defer cleanupPointUser(ctx, t, h3)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)
		u3 := seedPointUser(ctx, t, h3)
		defer cleanupOrderCreate(ctx, u1)

		// 第二个真实 SKU（同 SPU, 用于"SKU 不符"反例; 清理由 fixture 的 TF-SKU% 判据覆盖）
		sku2, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
			"sku_no": "TF-SKU-002", "spu_id": f.SpuId, "name": "TF-商品 白",
			"specs": `{"颜色":"白"}`, "price": "10.00", "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)

		const name = "TF-砍价成交下单"
		defer cleanupBargain(ctx, name)
		// 100/60/3 刀 → 3 刀到底价（86.67 → 73.34 → 60.00）
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 3)
		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		_, err = NewBargainLogic().Cut(ctx, u2, out.RecordId)
		t.AssertNil(err)

		addrId := seedUserAddress(ctx, t, u1)
		// 到底价前下单 → 40004
		_, err = NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
			RequestToken: "T-BGO-PRE", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, BargainRecordId: out.RecordId,
		})
		t.Assert(errCode(err), errcode.CodeBargainUnpayable)

		_, err = NewBargainLogic().Cut(ctx, u3, out.RecordId)
		t.AssertNil(err)
		p, err := NewBargainLogic().Progress(ctx, out.RecordId)
		t.AssertNil(err)
		t.Assert(p.Status, 2)

		// 数量不符（>1）→ 40004; SKU 与砍价单不符 → 40004
		_, err = NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
			RequestToken: "T-BGO-Q2", AddressId: addrId, SkuId: f.SkuId, Quantity: 2, BargainRecordId: out.RecordId,
		})
		t.Assert(errCode(err), errcode.CodeBargainUnpayable)
		_, err = NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
			RequestToken: "T-BGO-SKU", AddressId: addrId, SkuId: sku2, Quantity: 1, BargainRecordId: out.RecordId,
		})
		t.Assert(errCode(err), errcode.CodeBargainUnpayable)

		// 到底价成交: 按成交价 60.00 计价
		created, err := NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
			RequestToken: "T-BGO-OK", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, BargainRecordId: out.RecordId,
		})
		t.AssertNil(err)
		item, err := g.DB().GetOne(ctx, "SELECT price, pay_amount FROM trade_order_item WHERE order_no=?", created.OrderNo)
		t.AssertNil(err)
		t.Assert(item["price"].String(), "60.00") // 成交价 = 砍价到底价（非商品原价 10.00）
		t.Assert(item["pay_amount"].String(), "60.00")

		rec, err := g.DB().GetOne(ctx, "SELECT status, order_no FROM bargain_record WHERE id=?", out.RecordId)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 3)              // 已下单
		t.Assert(rec["order_no"].String(), created.OrderNo) // order_no 已回填

		// 重复下单（同砍价单再成交）→ 40004
		_, err = NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
			RequestToken: "T-BGO-DUP", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, BargainRecordId: out.RecordId,
		})
		t.Assert(errCode(err), errcode.CodeBargainUnpayable)
	})
}

// TestAssistLaunchOnce I1 回归: 并发双开发起 → 恰一次成功, 失败必须是 50006（唯一键兜底,
// 不得裸抛 10002）; per_limit 语义已被 uk_activity_user 钉死为"每人一次"。
func TestAssistLaunchOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ASL-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)

		const name = "TF-助力并发发起"
		defer cleanupAssist(ctx, name)
		actId := seedAssist(ctx, t, name, 1, 1, 1)

		var wg sync.WaitGroup
		errs := make([]error, 2)
		start := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				_, errs[idx] = NewAssistLogic().Launch(ctx, uid, actId)
			}(i)
		}
		close(start)
		wg.Wait()

		success := 0
		for _, e := range errs {
			if e == nil {
				success++
			} else {
				t.Assert(errCode(e), errcode.CodeAssistUsedUp) // 不允许 10002/其他
			}
		}
		t.Assert(success, 1)
	})
}

// TestFullReductionScopeFilter I3 回归: SpuId 范围命中过滤（全场命中 / 非命中排除）
// + ScopeDesc/EndTime 出口 + 过滤后的分页 total。
func TestFullReductionScopeFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)

		const nameA, nameB = "TF-满减全场", "TF-满减指定"
		cleanupFullReduction(ctx, nameA)
		cleanupFullReduction(ctx, nameB)
		defer cleanupFullReduction(ctx, nameA)
		defer cleanupFullReduction(ctx, nameB)
		seedFullReduction(ctx, t, nameA, "全场", 0)
		seedFullReduction(ctx, t, nameB, "指定商品", f.SpuId+999999) // 指向别的"商品", 不命中本 SPU

		// 不过滤: 两个都出, 且范围摘要/截止时间有值（I3: 契约出口）
		all, err := NewMarketingLogic().PublicFullReductions(ctx, 0, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.Assert(all.Total, int64(2))
		descA := ""
		for _, it := range all.List {
			t.Assert(it.EndTime != "", true)
			if it.Name == nameA {
				descA = it.ScopeDesc
				t.Assert(len(it.Ladders), 1)
			}
		}
		t.Assert(descA, "全场")

		// 按 SpuId 过滤: 只有"全场"命中
		filtered, err := NewMarketingLogic().PublicFullReductions(ctx, f.SpuId, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.Assert(filtered.Total, int64(1))
		t.Assert(filtered.List[0].Name, nameA)

		// 分页: total 为过滤后语义, 页大小 1 → 每页 1 条
		paged, err := NewMarketingLogic().PublicFullReductions(ctx, f.SpuId, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(paged.Total, int64(1))
		t.Assert(len(paged.List), 1)
	})
}

// seedFullReduction 造满减活动 + 一档 + 范围行（scopeKind: "全场"/"指定商品"; targetId 仅后者有效）。
func seedFullReduction(ctx context.Context, t *gtest.T, name, scopeKind string, targetId int64) {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO promotion_activity(name,start_time,end_time,status,deleted) "+
			"VALUES(?,DATE_SUB(NOW(), INTERVAL 1 HOUR),DATE_ADD(NOW(), INTERVAL 1 HOUR),1,0)", name)
	t.AssertNil(err)
	actId, _ := res.LastInsertId()
	_, err = g.DB().Exec(ctx,
		"INSERT INTO promotion_activity_ladder(activity_id,threshold_amount,discount_amount) VALUES(?,100.00,10.00)", actId)
	t.AssertNil(err)
	if scopeKind == "全场" {
		_, err = g.DB().Exec(ctx,
			"INSERT INTO promotion_activity_scope(activity_id,scope_type,target_id) VALUES(?,1,NULL)", actId)
	} else {
		_, err = g.DB().Exec(ctx,
			"INSERT INTO promotion_activity_scope(activity_id,scope_type,target_id) VALUES(?,3,?)", actId, targetId)
	}
	t.AssertNil(err)
}

func cleanupFullReduction(ctx context.Context, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_scope WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_ladder WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity WHERE name=?", name)
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

// TestBargainOrderConcurrentDoublePlace 复审新发现-7 守卫: 同一砍价单**并发**下单必须恰一次成交。
// 串行重复下单由事务外 loadBargainForOrder 的 status 校验兜住; 并发双花由下单事务内
// 步骤5.5 的条件置位（status 2→3, 判行数）兜住——本用例钉住的正是后者。
func TestBargainOrderConcurrentDoublePlace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const h1, h2, h3 = "BGC2-1", "BGC2-2", "BGC2-3"
		defer cleanupPointUser(ctx, t, h1)
		defer cleanupPointUser(ctx, t, h2)
		defer cleanupPointUser(ctx, t, h3)
		u1 := seedPointUser(ctx, t, h1)
		u2 := seedPointUser(ctx, t, h2)
		u3 := seedPointUser(ctx, t, h3)
		defer cleanupOrderCreate(ctx, u1)

		const name = "TF-砍价并发下单"
		defer cleanupBargain(ctx, name)
		_, itemId := seedBargain(ctx, t, name, f.SpuId, f.SkuId, "100.00", "60.00", 3)
		out, err := NewBargainLogic().Launch(ctx, u1, itemId)
		t.AssertNil(err)
		_, err = NewBargainLogic().Cut(ctx, u2, out.RecordId)
		t.AssertNil(err)
		_, err = NewBargainLogic().Cut(ctx, u3, out.RecordId)
		t.AssertNil(err) // 到底价

		addrId := seedUserAddress(ctx, t, u1)
		var wg sync.WaitGroup
		errs := make([]error, 2)
		orderNos := make([]string, 2)
		start := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				created, e := NewOrderLogic().Create(ctx, u1, model.OrderCreateInput{
					RequestToken:    fmt.Sprintf("T-BGC2-%d", idx),
					AddressId:       addrId,
					SkuId:           f.SkuId,
					Quantity:        1,
					BargainRecordId: out.RecordId,
				})
				if e == nil {
					orderNos[idx] = created.OrderNo
				}
				errs[idx] = e
			}(i)
		}
		close(start)
		wg.Wait()

		success := 0
		for i := 0; i < 2; i++ {
			if errs[i] == nil {
				success++
			} else {
				t.Assert(errCode(errs[i]), errcode.CodeBargainUnpayable) // 失败方必须是 40004
			}
		}
		t.Assert(success, 1) // 恰一次成交

		rec, err := g.DB().GetOne(ctx, "SELECT order_no FROM bargain_record WHERE id=?", out.RecordId)
		t.AssertNil(err)
		winner := orderNos[0]
		if winner == "" {
			winner = orderNos[1]
		}
		t.Assert(rec["order_no"].String(), winner) // 回填的必是成交那张单
	})
}
