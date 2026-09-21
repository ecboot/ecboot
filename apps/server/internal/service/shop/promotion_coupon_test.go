// promotion_coupon_test.go 优惠计算的元/分维度回归（012 修复轮 I7 及其补漏）。
// 根因: 优惠域金额列一律 DECIMAL(10,2) 存**元**（项目金额口径）, 而结算账本内部一律 int64 **分**。
// 原实现两处都没做换算——① 门槛（元）与商品金额（分）直接比 → 100 倍过宽;
// ② 抵扣额原样把"元"当"分"返回 → 100 倍少抵。评审仅指出券门槛一处, 补漏后两处齐修。
package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestCalcFullReductionFen 满减命中: 门槛按分比较; 抵扣额元→分换算。
func TestCalcFullReductionFen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity WHERE name='TF-满减活动'")
		actId, err := g.DB().Model("promotion_activity").Ctx(ctx).Data(g.Map{
			"name": "TF-满减活动", "status": 1, "deleted": 0,
			"start_time": gtime.Now().AddDate(0, 0, -1),
			"end_time":   gtime.Now().AddDate(0, 0, 1),
		}).InsertAndGetId()
		t.AssertNil(err)
		// AUTO_INCREMENT 复用被删活动的 id, 故先清掉该 id 下的孤儿档位（uk_activity_threshold 会撞）
		_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_ladder WHERE activity_id=?", actId)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_ladder WHERE activity_id=?", actId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity WHERE id=?", actId)
		}()
		_, err = g.DB().Model("promotion_activity_ladder").Ctx(ctx).Data(g.Map{
			"activity_id": actId, "threshold_amount": "100.00", "discount_amount": "20.00",
		}).Insert()
		t.AssertNil(err)

		// 150 分 = 1.5 元: 远未达"满 100 元"门槛 → 0
		// （修复前门槛是"元 vs 分"直接比, 1.5 元的车会命中满 100 减 20）
		t.Assert(calcFullReductionFen(ctx, "330000", 150), int64(0))
		// 9999 分 = 99.99 元: 差 1 分未达门槛 → 0
		t.Assert(calcFullReductionFen(ctx, "330000", 9999), int64(0))
		// 20000 分 = 200 元: 命中满 100 减 20 → 抵扣额 20 元 = 2000 分
		// （修复前返回 20——把"元"当"分", 少抵 100 倍）
		t.Assert(calcFullReductionFen(ctx, "330000", 20000), int64(2000))
	})
}

// TestCalcCouponDiscountFen 券抵扣: 门槛/抵扣额同维度; 过期券不得抵扣。
func TestCalcCouponDiscountFen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const uid = 8801
		_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name='TF-满100减20'")
		couponId, err := g.DB().Model("coupon").Ctx(ctx).Data(g.Map{
			"name": "TF-满100减20", "type": 1, "threshold_amount": "100.00", "discount_amount": "20.00",
			"total_count": 0, "per_limit": 1, "valid_type": 1, "status": 1, "deleted": 0,
		}).InsertAndGetId()
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM user_coupon WHERE coupon_id=?", couponId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE id=?", couponId)
		}()

		ucId, err := g.DB().Model("user_coupon").Ctx(ctx).Data(g.Map{
			"user_id": uid, "coupon_id": couponId, "status": 1,
			"expire_time": gtime.Now().AddDate(0, 0, 7),
		}).InsertAndGetId()
		t.AssertNil(err)

		// 1 元车（100 分）: 未达"满 100 元" → 0（修复前命中并被截断成 100 分, 等于白送）
		t.Assert(calcCouponDiscountFen(ctx, uid, ucId, 100), int64(0))
		// 100 元车（10000 分）: 刚达门槛 → 抵扣 20 元 = 2000 分（修复前返回 20 分）
		t.Assert(calcCouponDiscountFen(ctx, uid, ucId, 10000), int64(2000))
		// 50 元车（5000 分）: 未达门槛 → 0
		t.Assert(calcCouponDiscountFen(ctx, uid, ucId, 5000), int64(0))
		// 他人券 → 0
		t.Assert(calcCouponDiscountFen(ctx, uid+1, ucId, 10000), int64(0))

		// 过期券: expire_time 已过 → 0（原实现 select 了 expire_time 却从不比较）
		_, err = g.DB().Exec(ctx, "UPDATE user_coupon SET expire_time=? WHERE id=?",
			gtime.Now().AddDate(0, 0, -1), ucId)
		t.AssertNil(err)
		t.Assert(calcCouponDiscountFen(ctx, uid, ucId, 10000), int64(0))

		// 已使用券 → 0
		_, err = g.DB().Exec(ctx, "UPDATE user_coupon SET status=2, expire_time=? WHERE id=?",
			gtime.Now().AddDate(0, 0, 7), ucId)
		t.AssertNil(err)
		t.Assert(calcCouponDiscountFen(ctx, uid, ucId, 10000), int64(0))
	})
}
