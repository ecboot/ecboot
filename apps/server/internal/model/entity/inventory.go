// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Inventory is the golang structure for table inventory.
type Inventory struct {
	SkuId     uint64      `json:"skuId"     orm:"sku_id"     ` // SKU ID(与product_sku 1:1,主键即关联)
	Total     uint        `json:"total"     orm:"total"      ` // 总库存=可售+已锁定
	Locked    uint        `json:"locked"    orm:"locked"     ` // 下单锁定(未支付),可售=total-locked
	WarnCount uint        `json:"warnCount" orm:"warn_count" ` // 库存预警阈值(低于触发后台提醒,预留)
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
