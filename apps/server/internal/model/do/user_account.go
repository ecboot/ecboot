// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAccount is the golang structure of table user_account for DAO operations like Where/Data.
type UserAccount struct {
	g.Meta    `orm:"table:user_account, do:true"`
	Id        any         // 账户ID
	UserId    any         // 用户ID
	Balance   any         // 可用余额(有符号,欠款为负,后续佣金入账抵扣;欠款不可用于余额消费)
	Frozen    any         // 冻结金额(提现审核/打款中)
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
