// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Inventory is the golang structure of table inventory for DAO operations like Where/Data.
type Inventory struct {
	g.Meta    `orm:"table:inventory, do:true"`
	SkuId     any         // SKU ID(与product_sku 1:1,主键即关联)
	Total     any         // 总库存=可售+已锁定
	Locked    any         // 下单锁定(未支付),可售=total-locked
	WarnCount any         // 库存预警阈值(低于触发后台提醒,预留)
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
