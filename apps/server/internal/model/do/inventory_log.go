// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLog is the golang structure of table inventory_log for DAO operations like Where/Data.
type InventoryLog struct {
	g.Meta      `orm:"table:inventory_log, do:true"`
	Id          any         // 流水ID
	SkuId       any         // SKU ID
	OrderNo     any         // 关联订单号(后台调整等无单操作为空)
	ChangeType  any         // 变动类型:1下单锁定 2支付核销 3取消释放 4超时释放 5后台调整 6售后回补
	Quantity    any         // 变动数量(绝对值,方向由change_type决定)
	TotalAfter  any         // 变动后total快照(对账用)
	LockedAfter any         // 变动后locked快照(对账用)
	Operator    any         // 操作者:system/user:{id}/admin:{id}
	Remark      any         // 备注
	CreatedAt   *gtime.Time // 创建时间
}
