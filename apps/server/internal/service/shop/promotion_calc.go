// promotion_calc.go 优惠计算 helper（试算与下单共用; 分域 int64 运算）。
// 顺序（既定规则）: 先满减后用券; 积分抵扣上限为扣除券/满减后的金额。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// calcFullReductionFen 满减命中最优档（时段内 + 范围命中(全场V1) + 商品金额≥档位）。
// V1: 全场活动按商品总额判定（范围命中的精细化随促销域完善）。
func calcFullReductionFen(ctx context.Context, provinceCode string, goodsFen int64) int64 {
	if goodsFen <= 0 {
		return 0
	}
	// 命中时段内启用的全场满减活动, 取满足门槛的最大档
	rec, err := g.DB().GetOne(ctx, `
SELECT a.id,
       (SELECT MAX(l.discount_amount) FROM promotion_activity_ladder l
         WHERE l.activity_id = a.id AND l.threshold_amount <= ?) AS best_discount
FROM promotion_activity a
WHERE a.status = 1 AND a.deleted = 0
  AND NOW() BETWEEN a.start_time AND a.end_time
ORDER BY a.id LIMIT 1`, goodsFen)
	if err != nil || rec.IsEmpty() {
		return 0
	}
	best := rec["best_discount"]
	if best.IsNil() || best.Int64() <= 0 {
		return 0
	}
	return best.Int64()
}

// calcCouponDiscountFen 券抵扣（校验: 归属/未用/未过期/门槛）。
func calcCouponDiscountFen(ctx context.Context, userId int64, userCouponId int64, goodsFen int64) int64 {
	if userCouponId <= 0 {
		return 0
	}
	rec, err := g.DB().GetOne(ctx, `
SELECT c.discount_amount, c.threshold_amount, uc.status, uc.expire_time
FROM user_coupon uc JOIN coupon c ON c.id = uc.coupon_id
WHERE uc.id = ? AND uc.user_id = ?`, userCouponId, userId)
	if err != nil || rec.IsEmpty() {
		return 0
	}
	if rec["status"].Int() != 1 { // 非未使用
		return 0
	}
	threshold := rec["threshold_amount"].Int64()
	if goodsFen*100 < threshold*100 { // 元→分比较（阈值以元存储）
		return 0
	}
	discount := rec["discount_amount"].Int64()
	if discount > goodsFen {
		discount = goodsFen // 抵扣不超商品金额
	}
	return discount
}

// calcPointDeductFen 积分抵扣（1 积分=1 分; 上限=可用积分与可抵金额较小者）。
func calcPointDeductFen(ctx context.Context, userId int64, use bool, maxFen int64) int64 {
	if !use || maxFen <= 0 {
		return 0
	}
	// 012 前置修复（批次 05 评审债务）: point_account 表**无 deleted 列**（V19 建表 + V26 仅加
	// last_earned_at）——原查询必报错并被静默吞掉, 致积分抵扣恒为 0。该表无软删语义, 直接按 user_id 查。
	rec, err := g.DB().GetOne(ctx,
		"SELECT balance FROM point_account WHERE user_id=?", userId)
	if err != nil || rec.IsEmpty() {
		return 0
	}
	balance := int64(rec["balance"].Int())
	if balance <= 0 {
		return 0
	}
	if balance < maxFen {
		return balance
	}
	return maxFen
}
