// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCoupon is the golang structure of table user_coupon for DAO operations like Where/Data.
type UserCoupon struct {
	g.Meta     `orm:"table:user_coupon, do:true"`
	Id         any         // 用户券ID(订单user_coupon_id引用此表)
	UserId     any         // 持券用户ID
	CouponId   any         // 优惠券模板ID
	Status     any         // 状态:1未使用 2已使用 3已过期 4已退回(订单取消,可再用)
	ExpireTime *gtime.Time // 过期时间(领取时按模板规则计算落库)
	OrderNo    any         // 核销订单号
	UsedTime   *gtime.Time // 核销时间
	CreatedAt  *gtime.Time // 领取时间
	UpdatedAt  *gtime.Time // 更新时间
}
