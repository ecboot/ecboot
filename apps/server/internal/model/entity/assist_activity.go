// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistActivity is the golang structure for table assist_activity.
type AssistActivity struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 助力活动ID
	Name          string      `json:"name"          orm:"name"           ` // 活动名称
	RewardType    int         `json:"rewardType"    orm:"reward_type"    ` // 奖励类型:1优惠券 2积分
	RewardRef     uint64      `json:"rewardRef"     orm:"reward_ref"     ` // 奖励载体ID(如coupon.id;积分为点数存config)
	RequiredCount uint        `json:"requiredCount" orm:"required_count" ` // 所需助力人数
	PerLimit      uint        `json:"perLimit"      orm:"per_limit"      ` // 每人可发起次数
	StartTime     *gtime.Time `json:"startTime"     orm:"start_time"     ` // 开始时间(含)
	EndTime       *gtime.Time `json:"endTime"       orm:"end_time"       ` // 结束时间(不含)
	Config        string      `json:"config"        orm:"config"         ` // 扩展参数(积分数/阶梯奖励等)
	Status        int         `json:"status"        orm:"status"         ` // 状态:1启用 0停用
	Deleted       int         `json:"deleted"       orm:"deleted"        ` // 软删除:0否 1是
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` // 更新时间
}
