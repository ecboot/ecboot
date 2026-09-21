// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductCategory is the golang structure for table product_category.
type ProductCategory struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 分类ID
	ParentId  uint64      `json:"parentId"  orm:"parent_id"  ` // 父分类ID,0为根(固定三级:1/2/3)
	Name      string      `json:"name"      orm:"name"       ` // 分类名称
	Icon      string      `json:"icon"      orm:"icon"       ` // 分类图标URL
	Level     int         `json:"level"     orm:"level"      ` // 层级:1一级 2二级 3三级
	Sort      int         `json:"sort"      orm:"sort"       ` // 同级排序,越小越靠前
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0禁用
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
