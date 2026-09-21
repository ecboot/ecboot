// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFootprint is the golang structure of table user_footprint for DAO operations like Where/Data.
type UserFootprint struct {
	g.Meta     `orm:"table:user_footprint, do:true"`
	Id         any         // 足迹ID
	UserId     any         // 用户ID
	SpuId      any         // SPU ID
	ViewCount  any         // 累计浏览次数(重复浏览累加)
	LastViewAt *gtime.Time // 最近浏览时间(重复浏览更新此列,清理任务扫描)
	CreatedAt  *gtime.Time // 首次浏览时间
}
