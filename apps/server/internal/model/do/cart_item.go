// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartItem is the golang structure of table cart_item for DAO operations like Where/Data.
type CartItem struct {
	g.Meta    `orm:"table:cart_item, do:true"`
	Id        any         // 购物车项ID
	UserId    any         // 用户ID
	SkuId     any         // SKU ID
	Quantity  any         // 数量(应用层限制上限,如99)
	Checked   any         // 结算勾选:0未选 1选中
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
