// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistActivity is the golang structure of table assist_activity for DAO operations like Where/Data.
type AssistActivity struct {
	g.Meta        `orm:"table:assist_activity, do:true"`
	Id            any         // 助力活动ID
	Name          any         // 活动名称
	RewardType    any         // 奖励类型:1优惠券 2积分
	RewardRef     any         // 奖励载体ID(如coupon.id;积分为点数存config)
	RequiredCount any         // 所需助力人数
	PerLimit      any         // 每人可发起次数
	StartTime     *gtime.Time // 开始时间(含)
	EndTime       *gtime.Time // 结束时间(不含)
	Config        any         // 扩展参数(积分数/阶梯奖励等)
	Status        any         // 状态:1启用 0停用
	Deleted       any         // 软删除:0否 1是
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
