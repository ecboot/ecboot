// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCoupon is the golang structure for table user_coupon.
type UserCoupon struct {
	Id         uint64      `json:"id"         orm:"id"          ` // 用户券ID(订单user_coupon_id引用此表)
	UserId     uint64      `json:"userId"     orm:"user_id"     ` // 持券用户ID
	CouponId   uint64      `json:"couponId"   orm:"coupon_id"   ` // 优惠券模板ID
	Status     int         `json:"status"     orm:"status"      ` // 状态:1未使用 2已使用 3已过期 4已退回(订单取消,可再用)
	ExpireTime *gtime.Time `json:"expireTime" orm:"expire_time" ` // 过期时间(领取时按模板规则计算落库)
	OrderNo    string      `json:"orderNo"    orm:"order_no"    ` // 核销订单号
	UsedTime   *gtime.Time `json:"usedTime"   orm:"used_time"   ` // 核销时间
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 领取时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
}
