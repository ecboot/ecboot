// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InventoryLog is the golang structure for table inventory_log.
type InventoryLog struct {
	Id          uint64      `json:"id"          orm:"id"           ` // 流水ID
	SkuId       uint64      `json:"skuId"       orm:"sku_id"       ` // SKU ID
	OrderNo     string      `json:"orderNo"     orm:"order_no"     ` // 关联订单号(后台调整等无单操作为空)
	ChangeType  int         `json:"changeType"  orm:"change_type"  ` // 变动类型:1下单锁定 2支付核销 3取消释放 4超时释放 5后台调整 6售后回补
	Quantity    uint        `json:"quantity"    orm:"quantity"     ` // 变动数量(绝对值,方向由change_type决定)
	TotalAfter  uint        `json:"totalAfter"  orm:"total_after"  ` // 变动后total快照(对账用)
	LockedAfter uint        `json:"lockedAfter" orm:"locked_after" ` // 变动后locked快照(对账用)
	Operator    string      `json:"operator"    orm:"operator"     ` // 操作者:system/user:{id}/admin:{id}
	Remark      string      `json:"remark"      orm:"remark"       ` // 备注
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // 创建时间
}
