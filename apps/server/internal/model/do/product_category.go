// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductCategory is the golang structure of table product_category for DAO operations like Where/Data.
type ProductCategory struct {
	g.Meta    `orm:"table:product_category, do:true"`
	Id        any         // 分类ID
	ParentId  any         // 父分类ID,0为根(固定三级:1/2/3)
	Name      any         // 分类名称
	Icon      any         // 分类图标URL
	Level     any         // 层级:1一级 2二级 3三级
	Sort      any         // 同级排序,越小越靠前
	Status    any         // 状态:1启用 0禁用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
