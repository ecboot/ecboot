// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FlashSaleItem is the golang structure for table flash_sale_item.
type FlashSaleItem struct {
	Id         uint64      `json:"id"         orm:"id"          ` // 场次商品ID
	ActivityId uint64      `json:"activityId" orm:"activity_id" ` // 活动ID
	SkuId      uint64      `json:"skuId"      orm:"sku_id"      ` // SKU ID
	FlashPrice float64     `json:"flashPrice" orm:"flash_price" ` // 秒杀价(快照进订单项price)
	StockCount uint        `json:"stockCount" orm:"stock_count" ` // 活动限量(与inventory分账,契约R4)
	SoldCount  uint        `json:"soldCount"  orm:"sold_count"  ` // 已售(取消回补本列;可售=stock_count-sold_count)
	PerLimit   uint        `json:"perLimit"   orm:"per_limit"   ` // 每人限购
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
}
