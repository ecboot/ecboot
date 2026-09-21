// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointAccount is the golang structure for table point_account.
type PointAccount struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 积分账户ID
	UserId       uint64      `json:"userId"       orm:"user_id"        ` // 用户ID
	Balance      int         `json:"balance"      orm:"balance"        ` // 积分余额(有符号,回退欠款为负,禁止UNSIGNED)
	LastEarnedAt *gtime.Time `json:"lastEarnedAt" orm:"last_earned_at" ` // 最后获得时间(滚动有效期口径:自最后获得日起12个月内有效;过期任务清零并记流水)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     ` // 更新时间
}
