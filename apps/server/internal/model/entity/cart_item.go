// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItem is the golang structure for table cart_item.
type CartItem struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 购物车项ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 用户ID
	SkuId     uint64      `json:"skuId"     orm:"sku_id"     ` // SKU ID
	Quantity  uint        `json:"quantity"  orm:"quantity"   ` // 数量(应用层限制上限,如99)
	Checked   int         `json:"checked"   orm:"checked"    ` // 结算勾选:0未选 1选中
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
