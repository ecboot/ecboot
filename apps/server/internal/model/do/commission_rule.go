// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommissionRule is the golang structure of table commission_rule for DAO operations like Where/Data.
type CommissionRule struct {
	g.Meta     `orm:"table:commission_rule, do:true"`
	Id         any         // 规则ID
	ScopeType  any         // 作用域:1分类(默认) 2商品(覆盖)
	ScopeId    any         // 作用域目标ID(分类ID或SPU ID)
	Level1Rate any         // 一级佣金比例%(直接邀请人,0-100)
	Level2Rate any         // 二级佣金比例%(邀请人的邀请人,0-100)
	Status     any         // 状态:1启用 0停用
	Deleted    any         // 软删除:0否 1是
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}
