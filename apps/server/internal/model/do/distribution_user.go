// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DistributionUser is the golang structure of table distribution_user for DAO operations like Where/Data.
type DistributionUser struct {
	g.Meta    `orm:"table:distribution_user, do:true"`
	Id        any         // 推广员ID
	UserId    any         // 用户ID
	Status    any         // 状态:1待审核 2通过 3冻结(冻结不产生新佣金,存量可提现)
	Level     any         // 推广员等级(V1单一等级=1;多级启用阈值见system_config: distribution.level.threshold)
	ApplyTime *gtime.Time // 申请时间
	AuditTime *gtime.Time // 审核时间
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
