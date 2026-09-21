// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InviteRecord is the golang structure for table invite_record.
type InviteRecord struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 激励记录ID
	NewUserId     uint64      `json:"newUserId"     orm:"new_user_id"    ` // 新用户ID(仅可被激励一次)
	InviterId     uint64      `json:"inviterId"     orm:"inviter_id"     ` // 邀请人用户ID
	RewardType    int         `json:"rewardType"    orm:"reward_type"    ` // 奖励类型:1优惠券
	RewardTrigger int         `json:"rewardTrigger" orm:"reward_trigger" ` // 奖励触发时机:1注册即发 2首单后发(防刷,需订单回调触发)
	RewardRef     uint64      `json:"rewardRef"     orm:"reward_ref"     ` // 奖励载体ID(如user_coupon.id)
	Status        int         `json:"status"        orm:"status"         ` // 状态:1已发放(发放与订单同事务,幂等)
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
}
