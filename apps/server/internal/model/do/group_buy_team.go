// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyTeam is the golang structure of table group_buy_team for DAO operations like Where/Data.
type GroupBuyTeam struct {
	g.Meta       `orm:"table:group_buy_team, do:true"`
	Id           any         // 团ID
	ActivityId   any         // 活动ID
	LeaderUserId any         // 团长用户ID
	Status       any         // 状态:1拼团中 2已成团 3已解散(超时/取消)
	MemberCount  any         // 当前成员数(应用层与member表同事务维护)
	ExpireTime   *gtime.Time // 成团截止时间(超时扫描依据)
	SuccessTime  *gtime.Time // 成团时间
	CreatedAt    *gtime.Time // 开团时间
	UpdatedAt    *gtime.Time // 更新时间
}
