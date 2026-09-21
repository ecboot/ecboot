// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Coupon is the golang structure for table coupon.
type Coupon struct {
	Id              uint64      `json:"id"              orm:"id"               ` // 优惠券模板ID
	Name            string      `json:"name"            orm:"name"             ` // 券名称(如:满100减20)
	Type            int         `json:"type"            orm:"type"             ` // 券类型:1满减券 2无门槛券 3折扣券(预留)
	ThresholdAmount float64     `json:"thresholdAmount" orm:"threshold_amount" ` // 使用门槛(满X元可用,0=无门槛)
	DiscountAmount  float64     `json:"discountAmount"  orm:"discount_amount"  ` // 抵扣金额(满减/无门槛券)
	DiscountRate    float64     `json:"discountRate"    orm:"discount_rate"    ` // 折扣率0.01-0.99(折扣券预留,如0.85=85折)
	TotalCount      uint        `json:"totalCount"      orm:"total_count"      ` // 发放总量,0=不限量
	ReceivedCount   uint        `json:"receivedCount"   orm:"received_count"   ` // 已领取数量(领券事务内条件更新,防超发)
	PerLimit        uint        `json:"perLimit"        orm:"per_limit"        ` // 每人限领数量
	ValidType       int         `json:"validType"       orm:"valid_type"       ` // 有效期方式:1固定区间 2领取后N天
	ValidStartAt    *gtime.Time `json:"validStartAt"    orm:"valid_start_at"   ` // 固定区间开始(valid_type=1)
	ValidEndAt      *gtime.Time `json:"validEndAt"      orm:"valid_end_at"     ` // 固定区间结束(valid_type=1)
	ValidDays       uint        `json:"validDays"       orm:"valid_days"       ` // 领取后N天有效(valid_type=2)
	Status          int         `json:"status"          orm:"status"           ` // 状态:1启用 0停用(停用不影响已领取的券)
	Deleted         int         `json:"deleted"         orm:"deleted"          ` // 软删除:0否 1是
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 更新时间
}
