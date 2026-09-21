// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FlashSaleItem is the golang structure of table flash_sale_item for DAO operations like Where/Data.
type FlashSaleItem struct {
	g.Meta     `orm:"table:flash_sale_item, do:true"`
	Id         any         // 场次商品ID
	ActivityId any         // 活动ID
	SkuId      any         // SKU ID
	FlashPrice any         // 秒杀价(快照进订单项price)
	StockCount any         // 活动限量(与inventory分账,契约R4)
	SoldCount  any         // 已售(取消回补本列;可售=stock_count-sold_count)
	PerLimit   any         // 每人限购
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}
