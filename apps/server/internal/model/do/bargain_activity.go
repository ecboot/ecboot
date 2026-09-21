// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainActivity is the golang structure of table bargain_activity for DAO operations like Where/Data.
type BargainActivity struct {
	g.Meta    `orm:"table:bargain_activity, do:true"`
	Id        any         // 砍价活动ID
	Name      any         // 活动名称
	SpuId     any         // SPU ID
	StartTime *gtime.Time // 开始时间(含)
	EndTime   *gtime.Time // 结束时间(不含)
	Status    any         // 状态:1启用 0停用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
