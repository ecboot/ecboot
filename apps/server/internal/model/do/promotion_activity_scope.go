// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionActivityScope is the golang structure of table promotion_activity_scope for DAO operations like Where/Data.
type PromotionActivityScope struct {
	g.Meta       `orm:"table:promotion_activity_scope, do:true"`
	Id           any         // 范围ID
	ActivityId   any         // 活动ID
	ScopeType    any         // 范围类型:1全场 2分类 3商品
	TargetId     any         // 目标ID(scope_type=2分类ID/3商品SPU_ID;全场为NULL,每活动至多一条全场行)
	TargetIdNorm any         // 唯一键载体:全场(NULL)归一为0,堵住NULL≠NULL多行漏洞
	CreatedAt    *gtime.Time // 创建时间
}
