// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionActivityLadder is the golang structure of table promotion_activity_ladder for DAO operations like Where/Data.
type PromotionActivityLadder struct {
	g.Meta          `orm:"table:promotion_activity_ladder, do:true"`
	Id              any         // 档位ID
	ActivityId      any         // 活动ID
	ThresholdAmount any         // 门槛金额(满X元,范围内商品合计)
	DiscountAmount  any         // 优惠金额(减Y元)
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
