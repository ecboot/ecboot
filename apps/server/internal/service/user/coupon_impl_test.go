package user

import (
	"context"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 优惠券 fixture（012） ----

// seedCoupon 建券模板（validType 1固定区间/2领取后N天）, 返回 couponId。
func seedCoupon(ctx context.Context, t *gtest.T, name string, threshold, discount string,
	totalCount, perLimit, validType int, validDays int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO coupon(name,type,threshold_amount,discount_amount,total_count,received_count,per_limit,"+
			"valid_type,valid_days,status) VALUES(?,1,?,?,?,0,?,?,?,1)",
		name, threshold, discount, totalCount, perLimit, validType, validDays)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupCoupon(ctx context.Context, t *gtest.T, name string) {
	_, _ = g.DB().Exec(ctx,
		"DELETE FROM user_coupon WHERE coupon_id IN (SELECT id FROM coupon WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name=?", name)
}

// TestCouponAvailableAndReceive 可领列表与领取（FR-013）: 防超发+限领; 超限 → 50001。
func TestCouponAvailableAndReceive(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900008001"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "领券测试", 0)
		defer cleanupCoupon(ctx, t, "T-券A")

		cid := seedCoupon(ctx, t, "T-券A", "100.00", "10.00", 2, 1, 2, 30)
		t.AssertGT(cid, 0)

		// 可领列表（CanReceive=true）
		res, err := AvailableTemplates(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		found := false
		for _, it := range res.List {
			if it.CouponId == cid {
				found = true
				t.Assert(it.CanReceive, true)
				t.Assert(it.Discount, "10.00")
			}
		}
		t.Assert(found, true)

		// 领取 → 我的券新增未使用
		_, err = Receive(ctx, uid, cid)
		t.AssertNil(err)
		mine, err := Mine(ctx, uid, 1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		n := 0
		for _, it := range mine.List {
			if it.Name == "T-券A" {
				n++
				t.Assert(it.Status, 1)
				t.Assert(it.ExpireTime != "", true) // valid_type=2 领取后 30 天
			}
		}
		t.Assert(n, 1)

		// 再次领取 → 超个人限领（per_limit=1）→ 50001
		_, err = Receive(ctx, uid, cid)
		t.Assert(errCode(err), errcode.CodeCouponSoldOut)
	})
}

// TestCouponMineLazyExpire 我的券四态与惰性过期（FR-013）。
func TestCouponMineLazyExpire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900008002"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "过期测试", 0)
		defer cleanupCoupon(ctx, t, "T-券B")

		cid := seedCoupon(ctx, t, "T-券B", "50.00", "5.00", 10, 2, 2, 30)
		_, err := Receive(ctx, uid, cid)
		t.AssertNil(err)

		// 把过期时间改到昨天 → 查"未使用"不应出现该券；查"已过期"应出现（惰性判定）
		_, err = g.DB().Exec(ctx,
			"UPDATE user_coupon SET expire_time=DATE_SUB(NOW(), INTERVAL 1 DAY) WHERE coupon_id=? AND user_id=?",
			cid, uid)
		t.AssertNil(err)

		unused, err := Mine(ctx, uid, 1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		for _, it := range unused.List {
			t.AssertNE(it.Name, "T-券B") // 不在"未使用"
		}
		expired, err := Mine(ctx, uid, 3, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		hit := false
		for _, it := range expired.List {
			if it.Name == "T-券B" {
				hit = true
			}
		}
		t.Assert(hit, true) // 在"已过期"（惰性判定, 不改表）
	})
}

// TestCouponUsableConsumeReturn 匹配/核销/退回（FR-014）。
func TestCouponUsableConsumeReturn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900008003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "核销测试", 0)
		defer cleanupCoupon(ctx, t, "T-券小")
		defer cleanupCoupon(ctx, t, "T-券大")

		// 门槛 50 抵扣 5；门槛 30 抵扣 3（金额 40 时前者不满足门槛）
		small := seedCoupon(ctx, t, "T-券小", "30.00", "3.00", 10, 1, 2, 30)
		big := seedCoupon(ctx, t, "T-券大", "50.00", "5.00", 10, 1, 2, 30)
		_, err := Receive(ctx, uid, small)
		t.AssertNil(err)
		_, err = Receive(ctx, uid, big)
		t.AssertNil(err)

		// 商品金额 40: 仅门槛 30 的券可用（门槛 50 不满足）
		usable, err := UsableForOrder(ctx, uid, "40.00")
		t.AssertNil(err)
		t.Assert(len(usable), 1)
		t.Assert(usable[0].Name, "T-券小")

		// 商品金额 60: 两张都可用, 按抵扣降序（5 在前）
		usable, err = UsableForOrder(ctx, uid, "60.00")
		t.AssertNil(err)
		t.Assert(len(usable), 2)
		t.Assert(usable[0].Name, "T-券大")
		t.Assert(usable[1].Name, "T-券小")

		// 核销 → 已使用 + 绑单号
		t.AssertNil(Consume(ctx, uid, usable[0].UserCouponId, "ORD-COUPON-1"))
		ucId := usable[0].UserCouponId
		rec, err := g.DB().GetOne(ctx, "SELECT status, order_no FROM user_coupon WHERE id=?", ucId)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 2)
		t.Assert(rec["order_no"].String(), "ORD-COUPON-1")

		// 退回 → 未使用（有效期不变）
		before, err := g.DB().GetValue(ctx, "SELECT expire_time FROM user_coupon WHERE id=?", ucId)
		t.AssertNil(err)
		t.AssertNil(ReturnBack(ctx, ucId))
		rec, err = g.DB().GetOne(ctx, "SELECT status, order_no, expire_time FROM user_coupon WHERE id=?", ucId)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 1)
		t.Assert(rec["order_no"].String(), "")
		t.Assert(rec["expire_time"].String(), before.String()) // 有效期不变
	})
}

// TestCouponReceiveConcurrentLimit 并发领取限领（评审 I4 补测）: 并发下"个人限领"必须成立。
// 关键点: 模板 total_count=0（不限量）时**防超发的条件更新整段被跳过**, 于是限领校验是唯一闸门——
// 原实现用普通 SELECT 读已领数, RR 隔离下两个并发事务都读到 mine=0 都放行 → 同一用户领到 2 张。
// 修复: 先对模板行加悲观锁（LockUpdate）串行化同模板领取, 并以锁定读（FOR UPDATE）取已领数。
func TestCouponReceiveConcurrentLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900009001"
		defer cleanupMember(ctx, t, phone)
		defer cleanupCoupon(ctx, t, "T-券并发")
		uid := seedMember(ctx, t, phone, "并发领取", 0)

		// total_count=0 不限量 + per_limit=1
		cid := seedCoupon(ctx, t, "T-券并发", "0.00", "1.00", 0, 1, 2, 30)

		const n = 8
		var wg sync.WaitGroup
		errs := make([]error, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start // 尽量同时开跑
				_, errs[idx] = Receive(ctx, uid, cid)
			}(i)
		}
		close(start)
		wg.Wait()

		// 恰好 1 张: 其余全部 50001
		ok := 0
		for _, e := range errs {
			if e == nil {
				ok++
			} else {
				t.Assert(errCode(e), errcode.CodeCouponSoldOut)
			}
		}
		t.Assert(ok, 1)
		cnt, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM user_coupon WHERE user_id=? AND coupon_id=?", uid, cid)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1)
	})
}
