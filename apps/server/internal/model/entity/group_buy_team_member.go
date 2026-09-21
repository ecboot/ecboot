// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyTeamMember is the golang structure for table group_buy_team_member.
type GroupBuyTeamMember struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 团员ID
	TeamId    uint64      `json:"teamId"    orm:"team_id"    ` // 团ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 成员用户ID
	OrderNo   string      `json:"orderNo"   orm:"order_no"   ` // 成员订单号
	JoinTime  *gtime.Time `json:"joinTime"  orm:"join_time"  ` // 参团时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
}
