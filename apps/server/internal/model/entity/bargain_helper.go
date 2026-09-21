// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainHelper is the golang structure for table bargain_helper.
type BargainHelper struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 帮砍记录ID
	RecordId     uint64      `json:"recordId"     orm:"record_id"      ` // 砍价单ID
	HelperUserId uint64      `json:"helperUserId" orm:"helper_user_id" ` // 帮砍人用户ID
	CutAmount    float64     `json:"cutAmount"    orm:"cut_amount"     ` // 本刀砍掉金额(与record条件更新同事务)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 帮砍时间
}
