// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductBrand is the golang structure for table product_brand.
type ProductBrand struct {
	Id          uint64      `json:"id"          orm:"id"          ` // 品牌ID
	Name        string      `json:"name"        orm:"name"        ` // 品牌名称(存活期间唯一,软删后仍占用)
	Logo        string      `json:"logo"        orm:"logo"        ` // 品牌Logo URL
	Description string      `json:"description" orm:"description" ` // 品牌简介
	Sort        int         `json:"sort"        orm:"sort"        ` // 排序,越小越靠前
	Status      int         `json:"status"      orm:"status"      ` // 状态:1启用 0禁用
	Deleted     int         `json:"deleted"     orm:"deleted"     ` // 软删除:0否 1是
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` // 更新时间
}
