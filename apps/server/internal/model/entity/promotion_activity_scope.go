// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PromotionActivityScope is the golang structure for table promotion_activity_scope.
type PromotionActivityScope struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 范围ID
	ActivityId   uint64      `json:"activityId"   orm:"activity_id"    ` // 活动ID
	ScopeType    int         `json:"scopeType"    orm:"scope_type"     ` // 范围类型:1全场 2分类 3商品
	TargetId     uint64      `json:"targetId"     orm:"target_id"      ` // 目标ID(scope_type=2分类ID/3商品SPU_ID;全场为NULL,每活动至多一条全场行)
	TargetIdNorm int64       `json:"targetIdNorm" orm:"target_id_norm" ` // 唯一键载体:全场(NULL)归一为0,堵住NULL≠NULL多行漏洞
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 创建时间
}
