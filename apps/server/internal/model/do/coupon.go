// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Coupon is the golang structure of table coupon for DAO operations like Where/Data.
type Coupon struct {
	g.Meta          `orm:"table:coupon, do:true"`
	Id              any         // 优惠券模板ID
	Name            any         // 券名称(如:满100减20)
	Type            any         // 券类型:1满减券 2无门槛券 3折扣券(预留)
	ThresholdAmount any         // 使用门槛(满X元可用,0=无门槛)
	DiscountAmount  any         // 抵扣金额(满减/无门槛券)
	DiscountRate    any         // 折扣率0.01-0.99(折扣券预留,如0.85=85折)
	TotalCount      any         // 发放总量,0=不限量
	ReceivedCount   any         // 已领取数量(领券事务内条件更新,防超发)
	PerLimit        any         // 每人限领数量
	ValidType       any         // 有效期方式:1固定区间 2领取后N天
	ValidStartAt    *gtime.Time // 固定区间开始(valid_type=1)
	ValidEndAt      *gtime.Time // 固定区间结束(valid_type=1)
	ValidDays       any         // 领取后N天有效(valid_type=2)
	Status          any         // 状态:1启用 0停用(停用不影响已领取的券)
	Deleted         any         // 软删除:0否 1是
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
