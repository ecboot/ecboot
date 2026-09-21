// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommissionRule is the golang structure for table commission_rule.
type CommissionRule struct {
	Id         uint64      `json:"id"         orm:"id"          ` // 规则ID
	ScopeType  int         `json:"scopeType"  orm:"scope_type"  ` // 作用域:1分类(默认) 2商品(覆盖)
	ScopeId    uint64      `json:"scopeId"    orm:"scope_id"    ` // 作用域目标ID(分类ID或SPU ID)
	Level1Rate float64     `json:"level1Rate" orm:"level1_rate" ` // 一级佣金比例%(直接邀请人,0-100)
	Level2Rate float64     `json:"level2Rate" orm:"level2_rate" ` // 二级佣金比例%(邀请人的邀请人,0-100)
	Status     int         `json:"status"     orm:"status"      ` // 状态:1启用 0停用
	Deleted    int         `json:"deleted"    orm:"deleted"     ` // 软删除:0否 1是
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
}
