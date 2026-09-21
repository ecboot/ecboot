// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyItem is the golang structure for table group_buy_item.
type GroupBuyItem struct {
	Id         uint64      `json:"id"         orm:"id"          ` // 拼团场次商品ID
	ActivityId uint64      `json:"activityId" orm:"activity_id" ` // 拼团活动ID
	SkuId      uint64      `json:"skuId"      orm:"sku_id"      ` // SKU ID
	GroupPrice float64     `json:"groupPrice" orm:"group_price" ` // 该SKU成团价(下单快照进订单项price)
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
}
