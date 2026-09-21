// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FlashSaleActivity is the golang structure of table flash_sale_activity for DAO operations like Where/Data.
type FlashSaleActivity struct {
	g.Meta    `orm:"table:flash_sale_activity, do:true"`
	Id        any         // 活动ID
	Name      any         // 活动名称
	StartTime *gtime.Time // 开始时间(含)
	EndTime   *gtime.Time // 结束时间(不含)
	Status    any         // 状态:1启用 0停用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
