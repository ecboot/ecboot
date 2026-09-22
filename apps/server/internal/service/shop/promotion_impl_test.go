// promotion_impl_test.go 营销后台——券模板管理（016-marketing-admin 批次 10）。
// 重点: 类型/有效期条件校验矩阵、停发与软删对 C 端"可领券"的即时联动、记录下钻。
package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- fixture（自清判据=精确名; 子表先删） ----

func cleanupCouponFixture(ctx context.Context, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE uc FROM user_coupon uc JOIN coupon c ON uc.coupon_id=c.id WHERE c.name=?", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name=?", name)
}

// seedCoupon 直插一张进行中的券模板（供 Update/Records/联动用例复用）。
func seedCoupon(ctx context.Context, t *gtest.T, name string, status int) int64 {
	_, _ = g.DB().Exec(ctx,
		"INSERT INTO coupon(name,type,threshold_amount,discount_amount,total_count,received_count,per_limit,"+
			"valid_type,valid_start_at,valid_end_at,status,deleted) "+
			"VALUES(?,1,100.00,20.00,100,0,1,1,'2026-01-01 00:00:00','2026-12-31 23:59:59',?,0)", name, status)
	id, err := g.DB().GetValue(ctx, "SELECT id FROM coupon WHERE name=?", name)
	t.AssertNil(err)
	return id.Int64()
}

// ---- US1 券模板管理 ----

// TestAdminCouponCreateMatrix 创建校验矩阵（FR-1）。
func TestAdminCouponCreateMatrix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "TF-券矩阵"
		defer cleanupCouponFixture(ctx, n)

		logic := NewCouponLogic()
		// 合法: 满减券
		id, err := logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 1, Threshold: "100.00", Discount: "20.00", TotalCount: 100, PerLimit: 1,
			ValidType: 1, ValidStartAt: "2026-01-01 00:00:00", ValidEndAt: "2026-12-31 23:59:59",
		})
		t.AssertNil(err)
		t.Assert(id > 0, true)

		// type=1 缺门槛 → 10001
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 1, Discount: "20.00", ValidType: 1,
			ValidStartAt: "2026-01-01 00:00:00", ValidEndAt: "2026-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 无门槛券合法（threshold 缺省 0）
		id2, err := logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 2, Discount: "5.00", TotalCount: 10,
			ValidType: 2, ValidDays: 7,
		})
		t.AssertNil(err)
		t.Assert(id2 > 0, true)

		// validType=2 缺 validDays → 10001
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 2, Discount: "5.00", ValidType: 2, ValidDays: 0,
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 固定区间 start>=end → 10001
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 1, Threshold: "100.00", Discount: "20.00",
			ValidType: 1, ValidStartAt: "2026-12-31 00:00:00", ValidEndAt: "2026-01-01 00:00:00",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 金额 >2 位小数 → 10001
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 1, Threshold: "100.00", Discount: "20.005",
			ValidType: 1, ValidStartAt: "2026-01-01 00:00:00", ValidEndAt: "2026-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 折扣券类型(3)未开放 → 10001
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n, Type: 3, Discount: "0.85",
			ValidType: 1, ValidStartAt: "2026-01-01 00:00:00", ValidEndAt: "2026-12-31 23:59:59",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}

// TestAdminCouponLifecycle 列表(含已领数)/详情/停发联动 C 端/软删/记录下钻（FR-1/SC-3）。
func TestAdminCouponLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "TF-券全链路"
		defer cleanupCouponFixture(ctx, n)
		cid := seedCoupon(ctx, t, n, 1)

		logic := NewCouponLogic()

		// 列表可见 + 详情回读（validType 出口, D3-②）
		list, err := logic.AdminList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		found := false
		for _, it := range list.List {
			if it.Name == n {
				found = true
				t.Assert(it.Received, 0)
				t.Assert(it.TotalCount, 100)
			}
		}
		t.Assert(found, true)

		detail, err := logic.AdminDetail(ctx, cid)
		t.AssertNil(err)
		t.Assert(detail.Name, n)
		t.Assert(detail.ValidType, 1)
		t.Assert(detail.ValidDesc != "", true)
		_, err = logic.AdminDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// C1 回归（评审修复）: validType=2 的券 valid_start_at/valid_end_at 恒 NULL——
		// 原实现对 GTime() 直接解引用, 列表/详情 panic。建一张并列表/详情, 必须正常。
		const n2 = "TF-券NULL日期"
		defer cleanupCouponFixture(ctx, n2)
		_, err = logic.AdminCreate(ctx, model.CouponInput{
			Name: n2, Type: 2, Discount: "5.00", ValidType: 2, ValidDays: 7,
		})
		t.AssertNil(err)
		list2, err := logic.AdminList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err) // 不 panic
		hit2 := false
		for _, it := range list2.List {
			if it.Name == n2 {
				hit2 = true
				t.Assert(it.ValidType, 2)
			}
		}
		t.Assert(hit2, true)
		id2v, err := g.DB().GetValue(ctx, "SELECT id FROM coupon WHERE name=?", n2)
		t.AssertNil(err)
		d2, err := logic.AdminDetail(ctx, id2v.Int64())
		t.AssertNil(err)
		t.Assert(d2.ValidType, 2)

		// C 端"可领券"含该券（用 PublicList 分页口径断言, 不受首页 3 条限额的环境挤占影响;
		// PublicList 与首页券块同源同口径——M7 补实现后正好互为验证）
		pub, err := logic.PublicList(ctx, model.PageReq{Page: 1, PageSize: 100})
		t.AssertNil(err)
		hit := false
		for _, c := range pub.List {
			if c.Name == n {
				hit = true
			}
		}
		t.Assert(hit, true)

		// 停发 → C 端即时消失（显式 status=0; I6 三态语义: nil 不动, 0 才停）
		zero := 0
		t.AssertNil(logic.AdminUpdate(ctx, cid, model.CouponInput{Status: &zero}))
		pub, err = logic.PublicList(ctx, model.PageReq{Page: 1, PageSize: 100})
		t.AssertNil(err)
		hit = false
		for _, c := range pub.List {
			if c.Name == n {
				hit = true
			}
		}
		t.Assert(hit, false)

		// 造一条领取记录 → 记录下钻可见（停发不影响既有领取）
		_, err = g.DB().Exec(ctx,
			"INSERT INTO user_coupon(user_id,coupon_id,status,expire_time) VALUES(999001,?,1,'2026-12-31 00:00:00')", cid)
		t.AssertNil(err)
		recs, err := logic.AdminRecords(ctx, cid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(recs.Total, int64(1))
		t.Assert(recs.List[0].UserId, int64(999001))
		t.Assert(recs.List[0].Status, 1)

		// 软删: 券与记录查询消失、user_coupon 行保留
		t.AssertNil(logic.AdminDelete(ctx, cid))
		cnt, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM coupon WHERE name=? AND deleted=0", n)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 0)
		keep, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_coupon uc JOIN coupon c ON uc.coupon_id=c.id WHERE c.name=?", n)
		t.AssertNil(err)
		t.Assert(keep.Int(), 1) // 领取行未级联删

		// I6 回归（独立段）: 部分更新（只改名, 不传 status）不得静默停发
		const n3 = "TF-券部分更新"
		defer cleanupCouponFixture(ctx, n3)
		cid3 := seedCoupon(ctx, t, n3, 1)
		t.AssertNil(logic.AdminUpdate(ctx, cid3, model.CouponInput{Name: n3 + "-改名"}))
		d3, err := logic.AdminDetail(ctx, cid3)
		t.AssertNil(err)
		t.Assert(d3.Status, 1)     // 启停未被触碰（原实现无条件写 status → 静默停发）
		t.Assert(d3.Name, n3+"-改名")
	})
}
