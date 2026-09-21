// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyActivity is the golang structure of table group_buy_activity for DAO operations like Where/Data.
type GroupBuyActivity struct {
	g.Meta       `orm:"table:group_buy_activity, do:true"`
	Id           any         // 活动ID
	SpuId        any         // SPU ID(成团价对全SKU生效)
	Name         any         // 活动名称
	GroupSize    any         // 成团人数
	ValidStartAt *gtime.Time // 活动开始时间
	ValidEndAt   *gtime.Time // 活动结束时间
	PerLimit     any         // 每人限购(整个活动期)
	Status       any         // 状态:1启用 0停用
	Deleted      any         // 软删除:0否 1是
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
