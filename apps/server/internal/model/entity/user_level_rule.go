// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserLevelRule is the golang structure for table user_level_rule.
type UserLevelRule struct {
	Id              uint64      `json:"id"              orm:"id"               ` // 等级规则ID
	Name            string      `json:"name"            orm:"name"             ` // 等级名称(如 普通会员/白银/黄金)
	GrowthThreshold uint        `json:"growthThreshold" orm:"growth_threshold" ` // 成长值门槛(唯一,等级判定=满足的最高门槛)
	Benefits        string      `json:"benefits"        orm:"benefits"         ` // 等级权益配置(折扣/优先客服等,键值对)
	Status          int         `json:"status"          orm:"status"           ` // 状态:1启用 0停用
	Deleted         int         `json:"deleted"         orm:"deleted"          ` // 软删除:0否 1是
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 更新时间
}
