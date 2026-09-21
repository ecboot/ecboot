// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFavorite is the golang structure of table user_favorite for DAO operations like Where/Data.
type UserFavorite struct {
	g.Meta    `orm:"table:user_favorite, do:true"`
	Id        any         // 收藏ID
	UserId    any         // 用户ID
	SpuId     any         // SPU ID(实时联查现价与可售状态,不快照)
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 收藏时间
	UpdatedAt *gtime.Time // 更新时间
}
