// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainItem is the golang structure for table bargain_item.
type BargainItem struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 砍价场次商品ID
	ActivityId    uint64      `json:"activityId"    orm:"activity_id"    ` // 活动ID
	SkuId         uint64      `json:"skuId"         orm:"sku_id"         ` // SKU ID
	OriginalPrice float64     `json:"originalPrice" orm:"original_price" ` // 起始价(=发起时售价口径)
	FloorPrice    float64     `json:"floorPrice"    orm:"floor_price"    ` // 底价(砍到底价即可下单;floor<=original应用校验)
	MaxCutCount   uint        `json:"maxCutCount"   orm:"max_cut_count"  ` // 最大帮砍刀数(0=不限,金额收敛到底价)
	Config        string      `json:"config"        orm:"config"         ` // 玩法参数(随机递减算法边界/首刀上限等,应用层约定)
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` // 更新时间
}
