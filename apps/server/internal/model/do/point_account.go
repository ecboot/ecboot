// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointAccount is the golang structure of table point_account for DAO operations like Where/Data.
type PointAccount struct {
	g.Meta       `orm:"table:point_account, do:true"`
	Id           any         // 积分账户ID
	UserId       any         // 用户ID
	Balance      any         // 积分余额(有符号,回退欠款为负,禁止UNSIGNED)
	LastEarnedAt *gtime.Time // 最后获得时间(滚动有效期口径:自最后获得日起12个月内有效;过期任务清零并记流水)
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
