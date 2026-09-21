// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure for table user_account.
type UserAccount struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 账户ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 用户ID
	Balance   float64     `json:"balance"   orm:"balance"    ` // 可用余额(有符号,欠款为负,后续佣金入账抵扣;欠款不可用于余额消费)
	Frozen    float64     `json:"frozen"    orm:"frozen"     ` // 冻结金额(提现审核/打款中)
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
