// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistRecord is the golang structure of table assist_record for DAO operations like Where/Data.
type AssistRecord struct {
	g.Meta      `orm:"table:assist_record, do:true"`
	Id          any         // 助力参与记录ID
	ActivityId  any         // 活动ID
	UserId      any         // 发起人用户ID
	HelperCount any         // 已助力人数(达required_count触发发奖)
	Status      any         // 状态:1进行中 2已完成发奖 3已过期
	FinishTime  *gtime.Time // 完成时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
