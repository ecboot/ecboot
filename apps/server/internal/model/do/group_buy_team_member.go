// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyTeamMember is the golang structure of table group_buy_team_member for DAO operations like Where/Data.
type GroupBuyTeamMember struct {
	g.Meta    `orm:"table:group_buy_team_member, do:true"`
	Id        any         // 团员ID
	TeamId    any         // 团ID
	UserId    any         // 成员用户ID
	OrderNo   any         // 成员订单号
	JoinTime  *gtime.Time // 参团时间
	CreatedAt *gtime.Time // 创建时间
}
