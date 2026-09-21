// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistHelper is the golang structure of table assist_helper for DAO operations like Where/Data.
type AssistHelper struct {
	g.Meta       `orm:"table:assist_helper, do:true"`
	Id           any         // 助力记录ID
	RecordId     any         // 参与记录ID
	HelperUserId any         // 助力人用户ID
	CreatedAt    *gtime.Time // 助力时间
}
