// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserLevelRule is the golang structure of table user_level_rule for DAO operations like Where/Data.
type UserLevelRule struct {
	g.Meta          `orm:"table:user_level_rule, do:true"`
	Id              any         // 等级规则ID
	Name            any         // 等级名称(如 普通会员/白银/黄金)
	GrowthThreshold any         // 成长值门槛(唯一,等级判定=满足的最高门槛)
	Benefits        any         // 等级权益配置(折扣/优先客服等,键值对)
	Status          any         // 状态:1启用 0停用
	Deleted         any         // 软删除:0否 1是
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
