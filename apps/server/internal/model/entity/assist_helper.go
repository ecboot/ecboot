// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AssistHelper is the golang structure for table assist_helper.
type AssistHelper struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 助力记录ID
	RecordId     uint64      `json:"recordId"     orm:"record_id"      ` // 参与记录ID
	HelperUserId uint64      `json:"helperUserId" orm:"helper_user_id" ` // 助力人用户ID
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 助力时间
}
