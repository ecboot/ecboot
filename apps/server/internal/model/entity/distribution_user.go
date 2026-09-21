// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DistributionUser is the golang structure for table distribution_user.
type DistributionUser struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 推广员ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 用户ID
	Status    int         `json:"status"    orm:"status"     ` // 状态:1待审核 2通过 3冻结(冻结不产生新佣金,存量可提现)
	Level     int         `json:"level"     orm:"level"      ` // 推广员等级(V1单一等级=1;多级启用阈值见system_config: distribution.level.threshold)
	ApplyTime *gtime.Time `json:"applyTime" orm:"apply_time" ` // 申请时间
	AuditTime *gtime.Time `json:"auditTime" orm:"audit_time" ` // 审核时间
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
