// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductBrand is the golang structure of table product_brand for DAO operations like Where/Data.
type ProductBrand struct {
	g.Meta      `orm:"table:product_brand, do:true"`
	Id          any         // 品牌ID
	Name        any         // 品牌名称(存活期间唯一,软删后仍占用)
	Logo        any         // 品牌Logo URL
	Description any         // 品牌简介
	Sort        any         // 排序,越小越靠前
	Status      any         // 状态:1启用 0禁用
	Deleted     any         // 软删除:0否 1是
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
