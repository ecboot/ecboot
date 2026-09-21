// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistRecord is the golang structure for table assist_record.
type AssistRecord struct {
	Id          uint64      `json:"id"          orm:"id"           ` // 助力参与记录ID
	ActivityId  uint64      `json:"activityId"  orm:"activity_id"  ` // 活动ID
	UserId      uint64      `json:"userId"      orm:"user_id"      ` // 发起人用户ID
	HelperCount uint        `json:"helperCount" orm:"helper_count" ` // 已助力人数(达required_count触发发奖)
	Status      int         `json:"status"      orm:"status"       ` // 状态:1进行中 2已完成发奖 3已过期
	FinishTime  *gtime.Time `json:"finishTime"  orm:"finish_time"  ` // 完成时间
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` // 更新时间
}
