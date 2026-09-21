// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainHelper is the golang structure of table bargain_helper for DAO operations like Where/Data.
type BargainHelper struct {
	g.Meta       `orm:"table:bargain_helper, do:true"`
	Id           any         // 帮砍记录ID
	RecordId     any         // 砍价单ID
	HelperUserId any         // 帮砍人用户ID
	CutAmount    any         // 本刀砍掉金额(与record条件更新同事务)
	CreatedAt    *gtime.Time // 帮砍时间
}
