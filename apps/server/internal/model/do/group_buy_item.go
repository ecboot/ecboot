// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyItem is the golang structure of table group_buy_item for DAO operations like Where/Data.
type GroupBuyItem struct {
	g.Meta     `orm:"table:group_buy_item, do:true"`
	Id         any         // 拼团场次商品ID
	ActivityId any         // 拼团活动ID
	SkuId      any         // SKU ID
	GroupPrice any         // 该SKU成团价(下单快照进订单项price)
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
}
