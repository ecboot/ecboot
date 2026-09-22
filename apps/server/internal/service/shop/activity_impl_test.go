// activity_impl_test.go 营销后台——满减/拼团/秒杀/砍价/助力活动管理（016-marketing-admin 批次 10）。
// 重点: 嵌套配置往返（提交==回读）、全量替换与已售保护、唯一键 1062 转业务码、
// per_limit 语义声明、C 端可见性联动（SC-3）。
package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- fixture ----

func cleanupActivityFixture(ctx context.Context, name string) {
	// 子表先删; 满减三表与四类活动表各自按名清理（精确匹配）
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_scope WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_ladder WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity WHERE name=?", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM group_buy_item WHERE activity_id IN (SELECT id FROM group_buy_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM group_buy_activity WHERE name=?", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_item WHERE activity_id IN (SELECT id FROM flash_sale_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_activity WHERE name=?", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_item WHERE activity_id IN (SELECT id FROM bargain_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_activity WHERE name=?", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_helper WHERE record_id IN (SELECT id FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?))", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_record WHERE activity_id IN (SELECT id FROM assist_activity WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM assist_activity WHERE name=?", name)
}

// ---- US2 满减 ----

// TestAdminFullReductionRoundTrip 档位+范围嵌套往返/重复门槛 50008/范围三型校验/全量替换（FR-2）。
func TestAdminFullReductionRoundTrip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const n = "TF-后台满减"
		defer cleanupActivityFixture(ctx, n)

		logic := NewActivityLogic()
		id, err := logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "10.00"}, {Threshold: "200.00", Discount: "30.00"}},
			Scopes:  []model.PromotionScope{{ScopeType: 3, TargetId: f.SpuId}},
		})
		t.AssertNil(err)
		t.Assert(id > 0, true)

		// 详情往返: 档位按门槛升序、范围回读
		d, err := logic.FullReductionDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(len(d.Ladders), 2)
		t.Assert(d.Ladders[0].Threshold, "100.00")
		t.Assert(len(d.Scopes), 1)
		t.Assert(d.Scopes[0].ScopeType, 3)
		t.Assert(d.Scopes[0].TargetId, f.SpuId)

		// 重复门槛 → 50008
		_, err = logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "1.00"}, {Threshold: "100.00", Discount: "2.00"}},
		})
		t.Assert(errCode(err), errcode.CodeLadderDup)

		// scope 三型校验: 分类/商品必须带 target; 全场不得带 target
		_, err = logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "1.00"}},
			Scopes:  []model.PromotionScope{{ScopeType: 2, TargetId: 0}},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "1.00"}},
			Scopes:  []model.PromotionScope{{ScopeType: 1, TargetId: f.SpuId}},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 时间窗 start>=end → 拒绝
		_, err = logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-12-31 00:00:00", EndTime: "2026-01-01 00:00:00",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "1.00"}},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 全量替换: 2 档 → 1 档, 范围换全场（不落行=全场口径）
		st := 1
		t.AssertNil(logic.FullReductionUpdate(ctx, id, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59", Status: &st,
			Ladders: []model.PromotionLadder{{Threshold: "300.00", Discount: "50.00"}},
			Scopes:  nil,
		}))
		// I5 回归: "减"不得超过"满"（满100减200 → 10001; 计价侧封顶属 012 挂账, 配置面拦住）
		_, err = logic.FullReductionCreate(ctx, model.PromotionActivityInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2026-12-31 23:59:59",
			Ladders: []model.PromotionLadder{{Threshold: "100.00", Discount: "200.00"}},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		d, err = logic.FullReductionDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(len(d.Ladders), 1)
		t.Assert(d.Ladders[0].Threshold, "300.00")
		t.Assert(len(d.Scopes), 0)

		// 软删 → C 端满减列表消失（SC-3）
		t.AssertNil(logic.FullReductionDelete(ctx, id))
		fr, err := NewMarketingLogic().PublicFullReductions(ctx, 0, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		for _, it := range fr.List {
			t.Assert(it.Name != n, true)
		}
	})
}

// ---- US3 拼团/秒杀/砍价 ----

// TestAdminFlashSaleItems 秒杀场次商品设置/唯一键转业务码/全量替换/已售保护（FR-3/FR-5）。
func TestAdminFlashSaleItems(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const n = "TF-后台秒杀"
		defer cleanupActivityFixture(ctx, n)
		const h1 = "ADM-FS-1"
		defer cleanupPointUser(ctx, t, h1)
		_ = seedPointUser(ctx, t, h1)

		logic := NewActivityLogic()
		actId, err := logic.FlashSaleCreate(ctx, model.ActivityTimeInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.AssertNil(err)

		// 设置 2 个场次商品（fixture sku + 第二真实 sku）
		sku2, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
			"sku_no": "TF-SKU-002", "spu_id": f.SpuId, "name": "TF-商品 白",
			"specs": `{"颜色":"白"}`, "price": "8.00", "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		t.AssertNil(logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: f.SkuId, FlashPrice: "5.00", StockCount: 10, PerLimit: 2},
			{SkuId: sku2, FlashPrice: "4.00", StockCount: 5},
		}))
		cnt, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM flash_sale_item WHERE activity_id=?", actId)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 2)
		// perLimit 缺省 1
		pl, err := g.DB().GetValue(ctx, "SELECT per_limit FROM flash_sale_item WHERE activity_id=? AND sku_id=?", actId, sku2)
		t.AssertNil(err)
		t.Assert(pl.Int(), 1)

		// 同次提交重复 SKU → 1062 转业务码（50002）
		err = logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: f.SkuId, FlashPrice: "5.00", StockCount: 10},
			{SkuId: f.SkuId, FlashPrice: "6.00", StockCount: 10},
		})
		t.Assert(errCode(err), errcode.CodeActivityInvalid)

		// 秒杀价非法（>2 位小数）→ 10001
		err = logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: f.SkuId, FlashPrice: "5.005", StockCount: 10},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 全量替换为 1 个（sku2 被移除; 均未售出 → 允许）
		t.AssertNil(logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: f.SkuId, FlashPrice: "5.00", StockCount: 10, PerLimit: 2},
		}))
		cnt, err = g.DB().GetValue(ctx, "SELECT COUNT(*) FROM flash_sale_item WHERE activity_id=?", actId)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1)

		// 已售保护 + C2 回归（评审修复）: 置 sold_count=1 模拟已售
		// ① 移除已售行 → 拒绝; ② **保留已售行**的替换 → 行 id 不变 + sold_count 守恒
		// （原全删重插实现把保留行也删了重建: 行 id 变、sold_count 清零、订单引用悬空——评审探针 P3）
		_, err = g.DB().Exec(ctx, "UPDATE flash_sale_item SET sold_count=1 WHERE activity_id=?", actId)
		t.AssertNil(err)
		oldRow, err := g.DB().GetOne(ctx, "SELECT id, sold_count FROM flash_sale_item WHERE activity_id=? AND sku_id=?", actId, f.SkuId)
		t.AssertNil(err)
		err = logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: sku2, FlashPrice: "4.00", StockCount: 5},
		})
		t.Assert(errCode(err), errcode.CodeActivityInvalid) // ① 已售行不得移除
		// ② 保留已售行（改价）+ 换入新 SKU → 成功
		t.AssertNil(logic.FlashSaleSetItems(ctx, actId, []model.ActivitySkuInput{
			{SkuId: f.SkuId, FlashPrice: "6.00", StockCount: 20, PerLimit: 1},
			{SkuId: sku2, FlashPrice: "4.00", StockCount: 5},
		}))
		keptRow, err := g.DB().GetOne(ctx, "SELECT id, sold_count, flash_price, stock_count FROM flash_sale_item WHERE activity_id=? AND sku_id=?", actId, f.SkuId)
		t.AssertNil(err)
		t.Assert(keptRow["id"].Int64(), oldRow["id"].Int64())             // 行 id 不变（原位更新, 引用不悬空）
		t.Assert(keptRow["sold_count"].Int(), oldRow["sold_count"].Int()) // sold_count 守恒
		t.Assert(keptRow["flash_price"].String(), "6.00")                 // 配置确实更新
		t.Assert(keptRow["stock_count"].Int(), 20)

		// I3 回归: **同值更新**不得误报"不存在"（affected=0 是无变化不是缺失）
		st := 1
		t.AssertNil(logic.FlashSaleUpdate(ctx, actId, model.ActivityTimeInput{
			Name: n, StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59", Status: &st,
		}))
		// I6 回归: Status 不传（nil）→ 启停不变
		t.AssertNil(logic.FlashSaleUpdate(ctx, actId, model.ActivityTimeInput{Name: n}))
		fs2, err := g.DB().GetValue(ctx, "SELECT status FROM flash_sale_activity WHERE id=?", actId)
		t.AssertNil(err)
		t.Assert(fs2.Int(), 1)

		// 软删 → C 端消失（SC-3）
		t.AssertNil(logic.FlashSaleDelete(ctx, actId))
		fs, err := NewMarketingLogic().PublicFlashSales(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range fs.List {
			t.Assert(it.Name != n, true)
		}
	})
}

// TestAdminGroupBuyAndBargain 拼团/砍价活动管理 + 场次商品校验（FR-3）。
func TestAdminGroupBuyAndBargain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const gn, bn = "TF-后台拼团", "TF-后台砍价"
		defer cleanupActivityFixture(ctx, gn)
		defer cleanupActivityFixture(ctx, bn)

		logic := NewActivityLogic()
		// 拼团: groupSize<2 由 api 校验; service 侧同样拒绝
		_, err := logic.GroupBuyCreate(ctx, model.GroupBuyInput{
			Name: gn, SpuId: f.SpuId, GroupSize: 1,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		gid, err := logic.GroupBuyCreate(ctx, model.GroupBuyInput{
			Name: gn, SpuId: f.SpuId, GroupSize: 3, PerLimit: 2,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.AssertNil(err)
		// 缺 groupPrice → 拒绝
		err = logic.GroupBuySetItems(ctx, gid, []model.ActivitySkuInput{{SkuId: f.SkuId}})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		t.AssertNil(logic.GroupBuySetItems(ctx, gid, []model.ActivitySkuInput{
			{SkuId: f.SkuId, GroupPrice: "7.00"},
		}))
		gp, err := g.DB().GetValue(ctx, "SELECT group_price FROM group_buy_item WHERE activity_id=?", gid)
		t.AssertNil(err)
		t.Assert(gp.String(), "7.00")

		// 砍价: floor>=original → 拒绝; 合法设置落库
		bid, err := logic.BargainCreate(ctx, model.BargainActivityInput{
			Name: bn, SpuId: f.SpuId,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.AssertNil(err)
		err = logic.BargainSetItems(ctx, bid, []model.ActivitySkuInput{
			{SkuId: f.SkuId, OriginalPrice: "60.00", FloorPrice: "100.00", MaxCutCount: 9},
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		t.AssertNil(logic.BargainSetItems(ctx, bid, []model.ActivitySkuInput{
			{SkuId: f.SkuId, OriginalPrice: "100.00", FloorPrice: "60.00", MaxCutCount: 9},
		}))
		bp, err := g.DB().GetOne(ctx, "SELECT original_price, floor_price, max_cut_count FROM bargain_item WHERE activity_id=?", bid)
		t.AssertNil(err)
		t.Assert(bp["original_price"].String(), "100.00")
		t.Assert(bp["max_cut_count"].Int(), 9)

		// 列表（status 筛选 + 分页）
		bl, err := logic.BargainList(ctx, 1, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		found := false
		for _, it := range bl.List {
			if it.Name == bn {
				found = true
				t.Assert(it.SpuId, f.SpuId)
			}
		}
		t.Assert(found, true)
	})
}

// ---- US4 助力 ----

// TestAdminAssistActivity 助力奖励两型落库/perLimit 如实存储（FR-4/D5/D6）。
func TestAdminAssistActivity(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		const an = "TF-后台助力"
		defer cleanupActivityFixture(ctx, an)

		// 造一张目标券（rewardType=1）
		couponId := seedCoupon(ctx, t, "TF-助力奖励券", 1)
		defer cleanupCouponFixture(ctx, "TF-助力奖励券")

		logic := NewActivityLogic()
		// rewardType=1 缺 rewardRef → 拒绝
		_, err := logic.AssistCreate(ctx, model.AssistActivityInput{
			Name: an, RewardType: 1, RewardRef: 0, RequiredCount: 3,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		// rewardRef 指向不存在的券 → 拒绝
		_, err = logic.AssistCreate(ctx, model.AssistActivityInput{
			Name: an, RewardType: 1, RewardRef: 999999999, RequiredCount: 3,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeActivityNotFound)

		// 券奖励: reward_ref 营券 ID
		aid, err := logic.AssistCreate(ctx, model.AssistActivityInput{
			Name: an, RewardType: 1, RewardRef: couponId, RequiredCount: 3,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.AssertNil(err)
		rr, err := g.DB().GetValue(ctx, "SELECT reward_ref FROM assist_activity WHERE id=?", aid)
		t.AssertNil(err)
		t.Assert(rr.Int64(), couponId)

		// 积分奖励: reward_ref=0 + config JSON 存 pointAmount（D5）
		aid2, err := logic.AssistCreate(ctx, model.AssistActivityInput{
			Name: an, RewardType: 2, PointAmount: 100, RequiredCount: 2, PerLimit: 3,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.AssertNil(err)
		row, err := g.DB().GetOne(ctx, "SELECT reward_ref, config, per_limit FROM assist_activity WHERE id=?", aid2)
		t.AssertNil(err)
		t.Assert(row["reward_ref"].Int64(), 0)
		// 积分数落 config JSON（D5）
		pl, err := g.DB().GetValue(ctx, "SELECT JSON_EXTRACT(config, '$.pointAmount') FROM assist_activity WHERE id=?", aid2)
		t.AssertNil(err)
		t.Assert(pl.Int(), 100)
		// perLimit>1 如实存储（D6: 不静默改值; 语义由契约声明）
		t.Assert(row["per_limit"].Int(), 3)

		// requiredCount<1 → 拒绝
		_, err = logic.AssistCreate(ctx, model.AssistActivityInput{
			Name: an, RewardType: 2, PointAmount: 100, RequiredCount: 0,
			StartTime: "2026-01-01 00:00:00", EndTime: "2099-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 列表 + C 端联动（rewardDesc 含积分信息, SC-3）
		al, err := logic.AssistList(ctx, 1, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		found := false
		for _, it := range al.List {
			if it.Name == an && it.RewardType == 2 {
				found = true
			}
		}
		t.Assert(found, true)
		pa, err := NewMarketingLogic().PublicAssists(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		hit := false
		for _, it := range pa.List {
			if it.Name == an {
				hit = true
			}
		}
		t.Assert(hit, true)

		// 软删 → C 端消失
		t.AssertNil(logic.AssistDelete(ctx, aid2))
		pa, err = NewMarketingLogic().PublicAssists(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range pa.List {
			t.Assert(it.Name != an || it.RewardType != 2, true)
		}
	})
}
