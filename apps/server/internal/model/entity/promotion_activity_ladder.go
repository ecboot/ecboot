// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionActivityLadder is the golang structure for table promotion_activity_ladder.
type PromotionActivityLadder struct {
	Id              uint64      `json:"id"              orm:"id"               ` // 档位ID
	ActivityId      uint64      `json:"activityId"      orm:"activity_id"      ` // 活动ID
	ThresholdAmount float64     `json:"thresholdAmount" orm:"threshold_amount" ` // 门槛金额(满X元,范围内商品合计)
	DiscountAmount  float64     `json:"discountAmount"  orm:"discount_amount"  ` // 优惠金额(减Y元)
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 更新时间
}
