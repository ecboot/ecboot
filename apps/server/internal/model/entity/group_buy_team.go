// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyTeam is the golang structure for table group_buy_team.
type GroupBuyTeam struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 团ID
	ActivityId   uint64      `json:"activityId"   orm:"activity_id"    ` // 活动ID
	LeaderUserId uint64      `json:"leaderUserId" orm:"leader_user_id" ` // 团长用户ID
	Status       int         `json:"status"       orm:"status"         ` // 状态:1拼团中 2已成团 3已解散(超时/取消)
	MemberCount  uint        `json:"memberCount"  orm:"member_count"   ` // 当前成员数(应用层与member表同事务维护)
	ExpireTime   *gtime.Time `json:"expireTime"   orm:"expire_time"    ` // 成团截止时间(超时扫描依据)
	SuccessTime  *gtime.Time `json:"successTime"  orm:"success_time"   ` // 成团时间
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 开团时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     ` // 更新时间
}
